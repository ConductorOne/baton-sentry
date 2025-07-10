package client

import (
	"context"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type Client struct {
	*uhttp.BaseHttpClient
}

func New(ctx context.Context, apiToken string) (*Client, error) {
	httpClient, err := uhttp.NewBearerAuth(apiToken).GetClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	return &Client{
		BaseHttpClient: uhttp.NewBaseHttpClient(httpClient),
	}, nil
}
