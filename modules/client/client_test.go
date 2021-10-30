package client_test

import (
	"context"
	"os"
	"testing"

	"github.com/HITSZ-PAM/pamcli/modules/client"
)

func Test_client_Resolve(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Valid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := client.Config{
				ServerAddr:   os.Getenv("PAM_SERVER_URL"),
				ClientID:     os.Getenv("PAM_CLIENT_ID"),
				ClientSecret: os.Getenv("PAM_CLIENT_SECRET"),
			}
			ctx := context.Background()
			c, err := client.NewClient(ctx, &cfg)
			if err != nil {
				t.Errorf("NewClient() error: %v", err)
				return
			}
			_, _, err = c.Resolve(os.Getenv("PAM_ACCOUNT_ID"))
			if err != nil {
				t.Errorf("client.Resolve() error: %v", err)
				return
			}
		})
	}
}

func Benchmark_client_Resolve(b *testing.B) {
	tests := []struct {
		name string
	}{
		{name: "Valid"},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			cfg := client.Config{
				ServerAddr:   os.Getenv("PAM_SERVER_URL"),
				ClientID:     os.Getenv("PAM_CLIENT_ID"),
				ClientSecret: os.Getenv("PAM_CLIENT_SECRET"),
			}
			ctx := context.Background()
			c, err := client.NewClient(ctx, &cfg)
			if err != nil {
				b.Errorf("NewClient() error: %v", err)
				return
			}
			_, _, err = c.Resolve(os.Getenv("PAM_ACCOUNT_ID"))
			if err != nil {
				b.Errorf("client.Resolve() error: %v", err)
				return
			}
			for i := 0; i < b.N; i++ {
				_, _, _ = c.Resolve(os.Getenv("PAM_ACCOUNT_ID"))
			}
		})
	}
}
