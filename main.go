package proxychecker

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
	"golang.org/x/xerrors"
)

var (
	Method string = http.MethodGet
	Target string = "https://httpbin.org/status/200"
)

type Type int

const (
	TypeHTTP Type = iota
	TypeSOCKS5
)

func Check(ctx context.Context, proxyAddr string, proxyType Type) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, Method, Target, nil)
	if err != nil {
		return false, xerrors.Errorf("failed to create request: %w", err)
	}

	var client *http.Client

	switch proxyType {
	default:
	case TypeHTTP:
		proxyURL, err := url.Parse(fmt.Sprintf("http://%s", proxyAddr))
		if err != nil {
			return false, xerrors.Errorf("failed to parse URL: %w", err)
		}

		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}

	case TypeSOCKS5:
		p, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
		if err != nil {
			return false, xerrors.Errorf("failed to initialize proxy: %w", err)
		}

		client = &http.Client{
			Transport: &http.Transport{
				Dial: p.Dial,
			},
		}
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
