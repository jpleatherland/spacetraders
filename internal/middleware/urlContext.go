package middleware

import (
	"context"
	"net/http"
)

const UrlKey = contextKey("url")

func UrlContext(url string, req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), UrlKey, url)
	return req.WithContext(ctx)
}

func GetUrlContext(ctx context.Context) (string, bool) {

	url, ok := ctx.Value(UrlKey).(string)
	return url, ok
}
