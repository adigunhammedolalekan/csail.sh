package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/saas/hostgolang/pkg/config"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	Set(key, value string) error
}

type defaultProxyClient struct {
	httpClient *http.Client
	cfg        *config.Config
}

func NewProxyClient(cfg *config.Config) (Client, error) {
	if _, err := url.Parse(cfg.ProxyServerAddress); err != nil {
		return nil, err
	}
	if len(cfg.ProxySecret) == 0 {
		return nil, errors.New("proxy secret is missing")
	}
	c := &http.Client{Timeout: 30 * time.Second}
	return &defaultProxyClient{httpClient: c, cfg: cfg}, nil
}

func (c *defaultProxyClient) Set(key, value string) error {
	if key == "" {
		return errors.New("name is empty")
	}
	// makes sure serviceURL is a valid http address/link
	if _, err := url.Parse(value); err != nil {
		return err
	}

	type payload struct {
		Name       string `json:"name"`
		ServiceUrl string `json:"service_url"`
	}
	p := &payload{Name: key, ServiceUrl: value}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		return err
	}
	u := fmt.Sprintf("%s/set", c.cfg.ProxyServerAddress)
	req, err := http.NewRequest("POST", u, buf)
	if err != nil {
		return err
	}
	req.Header.Add("X-Proxy-Secret", c.cfg.ProxySecret)
	r, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed to contact proxy front; returned http code %d", r.StatusCode))
	}
	return nil
}
