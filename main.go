package proxychecker

import (
	"context"
	"net/http"
	"net/url"

	"golang.org/x/xerrors"
)

var (
	Method string = http.MethodGet
	Target string = "https://httpbin.org/status/200"
)

func Check(ctx context.Context, proxy *url.URL) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, Method, Target, nil)
	if err != nil {
		return false, xerrors.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return false, xerrors.Errorf("failed to send request: %w", err)
	}

	if err := res.Body.Close(); err != nil {
		return false, xerrors.Errorf("failed to close response body: %w", err)
	}

	return true, nil
}
