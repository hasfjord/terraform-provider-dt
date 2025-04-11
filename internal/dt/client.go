// Copyright (c) HashiCorp, Inc.

package dt

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt/oidc"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	URL          string
	httpClient   http.Client
	oidc         *oidc.Client
	retryAfter   *retryAfter
	version      string
	rulesCache   *rulesCache
	projectCache *projectCache
}

type retryAfter struct {
	// retryAfter is the time after which we can send another request
	// after receiving a 429 Too Many Requests response
	t  time.Time
	mu sync.RWMutex
}

type Config struct {
	Oidc    oidc.Config
	URL     string
	Version string
}

func NewClient(cfg Config) *Client {
	return &Client{
		URL:        cfg.URL,
		httpClient: *http.DefaultClient,
		oidc:       oidc.NewClient(cfg.Oidc),
		retryAfter: &retryAfter{
			t:  time.Now(),
			mu: sync.RWMutex{},
		},
		version: cfg.Version,
		rulesCache: &rulesCache{
			notificationRules: make(map[string]NotificationRule),
			mu:                sync.RWMutex{},
		},
		projectCache: &projectCache{
			projects: make(map[string]Project),
			mu:       sync.RWMutex{},
		},
	}
}

func (c *Client) WithHttpClient(httpClient http.Client) *Client {
	c.httpClient = httpClient
	return c
}

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %d: %s", e.StatusCode, e.Body)
}

// time returns the retry after time.
func (r *retryAfter) time() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.t
}

// setTime sets the retry after time to the given time.
func (r *retryAfter) setTime(t time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.t = t
}

func (c *Client) DoRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	// Check if we need to wait for the retry after time
	// before sending the request
	time.Sleep(time.Until(c.retryAfter.time()))

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	ctx = tflog.SetField(ctx, "method", method)
	ctx = tflog.SetField(ctx, "url", url)

	tflog.Debug(ctx, "sending request to DT API")

	// Get an OIDC token and set it as a Bearer token in the request
	token, err := c.oidc.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to get OIDC token: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", fmt.Sprintf("TerraformProviderDT/%s(%s)", c.version, runtime.Version()))

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to send request: %w", err)
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to read response body: %w, status: %d", err, response.StatusCode)
	}
	ctx = tflog.SetField(ctx, "status_code", response.StatusCode)
	ctx = tflog.SetField(ctx, "body", string(bodyBytes))
	if response.StatusCode != http.StatusOK {
		tflog.Debug(ctx, "received non-200 status code from DT API")
		if response.StatusCode == http.StatusTooManyRequests {
			c.retryAfter.setTime(getRetryAfterTime(response))
			tflog.Debug(ctx, "received 429 status code from DT API, retrying request")

			// Retry the request
			return c.DoRequest(ctx, method, url, body)
		}
		return nil, &HTTPError{
			StatusCode: response.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	tflog.Debug(ctx, "received response from DT API")

	return bodyBytes, nil
}

// getRetryAfterTime gets the retry after time from a request by parsing the Retry-After header.
func getRetryAfterTime(res *http.Response) time.Time {
	retryAfter := res.Header.Get("Retry-After")
	if retryAfter == "" {
		return time.Now()
	}
	// Retry-After can be either a number of seconds or a date
	retryAfterDuration, err := time.ParseDuration(retryAfter + "s")
	if err != nil {
		retryAfterTime, err := http.ParseTime(retryAfter)
		if err != nil {
			return time.Now()
		}
		return retryAfterTime
	}
	return time.Now().Add(retryAfterDuration)
}
