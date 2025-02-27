package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt/oidc"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := realMain(ctx)
	done()
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	} else {
		fmt.Println("done")
	}
}

func realMain(ctx context.Context) error {
	cfg := dt.Config{
		URL: "https://api.dev.disruptive-technologies.com/v2",
		Oidc: oidc.Config{
			TokenEndpoint: "https://identity.dev.disruptive-technologies.com/oauth2/token",
			ClientID:      "ct5lk8324te000b24stg",
			ClientSecret:  os.Getenv("DT_OIDC_CLIENT_SECRET"),
			Email:         "ct5lk6j24te000b24teg@cka1u2lk2aecsarq62s0.serviceaccount.d21s.com",
		},
	}

	client := dt.NewClient(cfg)

	_, err := client.GetProject(ctx, "projects/ccol8iuk9smqiha4e8l0")
	if err != nil {
		return err
	}

	return nil
}
