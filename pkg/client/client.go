package client

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const DefaultBaseURL = "https://sentry.io/api/0/"

type sentryError struct {
	Detail json.RawMessage `json:"detail"`
	raw    json.RawMessage
}

func (e *sentryError) UnmarshalJSON(data []byte) error {
	e.raw = data
	type plain sentryError
	return json.Unmarshal(data, (*plain)(e))
}

func (e *sentryError) Message() string {
	if e.Detail != nil {
		var s string
		if err := json.Unmarshal(e.Detail, &s); err == nil {
			return s
		}
		var obj struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(e.Detail, &obj); err == nil && obj.Message != "" {
			return obj.Message
		}
		return string(e.Detail)
	}
	if len(e.raw) > 0 {
		return string(e.raw)
	}
	return ""
}

type Client struct {
	*uhttp.BaseHttpClient
	baseURL string
}

func New(ctx context.Context, apiToken string, baseURL string) (*Client, error) {
	httpClient, err := uhttp.NewBearerAuth(apiToken).GetClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return &Client{
		BaseHttpClient: uhttp.NewBaseHttpClient(httpClient),
		baseURL:        baseURL,
	}, nil
}

func (c *Client) doRequest(req *http.Request, target interface{}, options ...uhttp.DoOption) (*http.Response, *v2.RateLimitDescription, error) {
	var ratelimitData v2.RateLimitDescription
	doOptions := []uhttp.DoOption{
		uhttp.WithErrorResponse(&sentryError{}),
		uhttp.WithRatelimitData(&ratelimitData),
	}
	if target != nil {
		doOptions = append(doOptions, uhttp.WithJSONResponse(target))
	}
	doOptions = append(doOptions, options...)

	res, err := c.Do(req, doOptions...)
	if err != nil {
		return nil, &ratelimitData, err
	}
	defer res.Body.Close()
	return res, &ratelimitData, nil
}
