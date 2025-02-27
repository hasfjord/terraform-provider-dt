package dt

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt/oidc"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	URL        string
	httpClient http.Client
	oidc       *oidc.Client
}

type Config struct {
	Oidc oidc.Config
	URL  string
}

func NewClient(cfg Config) *Client {
	return &Client{
		URL:        cfg.URL,
		httpClient: *http.DefaultClient,
		oidc:       oidc.NewClient(cfg.Oidc),
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

func (c *Client) DoRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
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

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to send request: %w", err)
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to read response body: %w, status: %d", err, response.StatusCode)
	}
	if response.StatusCode != http.StatusOK {
		ctx = tflog.SetField(ctx, "status_code", response.StatusCode)
		ctx = tflog.SetField(ctx, "body", string(bodyBytes))
		tflog.Debug(ctx, "received non-200 status code from DT API")
		return nil, &HTTPError{
			StatusCode: response.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	ctx = tflog.SetField(ctx, "status_code", response.StatusCode)
	ctx = tflog.SetField(ctx, "body", string(bodyBytes))

	tflog.Debug(ctx, "received response from DT API")

	return bodyBytes, nil
}
