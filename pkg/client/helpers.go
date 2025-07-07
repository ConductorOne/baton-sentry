package client

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/peterhellberg/link"
	"go.uber.org/zap"
)

func logBody(ctx context.Context, bodyCloser io.ReadCloser) {
	l := ctxzap.Extract(ctx)
	body := make([]byte, 1024*1024)
	n, err := bodyCloser.Read(body)
	if err != nil {
		l.Error("error reading response body", zap.Error(err))
		return
	}
	l.Info("response body: ", zap.String("body", string(body[:n])))
}

// https://docs.sentry.io/api/pagination/
func HasNextPage(res *http.Response) bool {
	for _, l := range link.ParseResponse(res) {
		if v, ok := l.Extra["results"]; ok && v == "true" {
			return true
		}
	}
	return false
}

// https://docs.sentry.io/api/pagination/
func NextCursor(res *http.Response) string {
	for _, l := range link.ParseResponse(res) {
		if l.Rel == "next" {
			if v, ok := l.Extra["cursor"]; ok {
				return v
			}
		}
	}
	return ""
}
