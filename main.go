package proxychecker

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/xerrors"
	"h12.io/socks"
)

var (
	Method string = http.MethodGet
	Target string = "https://api.myip.com/"
)

type Type int

const (
	TypeHTTP Type = iota
	TypeSOCKS5
)

var client *http.Client

func Check(proxyAddr string, proxyType Type) (bool, error) {
	req, err := http.NewRequest(Method, Target, nil)
	if err != nil {
		return false, xerrors.Errorf("failed to create request: %w", err)
	}

	var transport *http.Transport

	switch proxyType {
	default:
	case TypeHTTP:
		proxyURL, err := url.Parse(fmt.Sprintf("http://%s", proxyAddr))
		if err != nil {
			return false, xerrors.Errorf("failed to parse URL: %w", err)
		}

		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

	case TypeSOCKS5:
		p := socks.Dial(fmt.Sprintf("socks5://%s", proxyAddr))
		transport = &http.Transport{
			Dial: p,
		}
	}

	transport.DisableKeepAlives = true
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	req.Header.Add("Connection", "close")

	client = &http.Client{
		Timeout:   time.Second * 20,
		Transport: transport,
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
