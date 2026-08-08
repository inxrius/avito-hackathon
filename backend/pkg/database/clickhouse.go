package database

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClickHouseHTTP struct {
	Endpoint string
	Database string
	Username string
	Password string
	Client   *http.Client
}

func NewClickHouseHTTP(endpoint, database, username, password string) *ClickHouseHTTP {
	return &ClickHouseHTTP{
		Endpoint: strings.TrimRight(endpoint, "/"),
		Database: database,
		Username: username,
		Password: password,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ClickHouseHTTP) Ping(ctx context.Context) error {
	_, err := c.Execute(ctx, "SELECT 1")
	return err
}

func (c *ClickHouseHTTP) Execute(ctx context.Context, query string) ([]byte, error) {
	if strings.TrimSpace(c.Endpoint) == "" {
		return nil, fmt.Errorf("clickhouse endpoint is empty")
	}
	endpoint, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse clickhouse endpoint: %w", err)
	}
	params := endpoint.Query()
	if c.Database != "" {
		params.Set("database", c.Database)
	}
	endpoint.RawQuery = params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewBufferString(query))
	if err != nil {
		return nil, fmt.Errorf("create clickhouse request: %w", err)
	}
	request.Header.Set("Content-Type", "text/plain; charset=utf-8")
	if c.Username != "" {
		request.SetBasicAuth(c.Username, c.Password)
	}

	client := c.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("execute clickhouse request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 16<<20))
	if err != nil {
		return nil, fmt.Errorf("read clickhouse response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("clickhouse status %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}
	return body, nil
}
