package requester

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Gate struct {
	client http.Client
}

func New() *Gate {
	return &Gate{}
}

func (g *Gate) Do(ctx context.Context, URL *url.URL) (str string, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", URL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("can't create http request: %s", err.Error())
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request error: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't read response body: %s", err.Error())
	}

	return string(body), err
}
