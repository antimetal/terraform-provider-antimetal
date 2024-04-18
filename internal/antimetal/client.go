// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package antimetal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultURL     = "https://groups.api.antimetal.com"
	defaultTimeout = 30 * time.Second
)

type ClientOption func(*Client) error

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		if c.httpClient == nil {
			return fmt.Errorf("Client is missing httpClient, probably because it wasn't properly initialized")
		}

		c.httpClient.Timeout = timeout
		return nil
	}
}

func WithURL(url string) ClientOption {
	return func(c *Client) error {
		c.url = url
		return nil
	}
}

type Client struct {
	url        string
	httpClient *http.Client
}

func NewClient(opts ...ClientOption) (*Client, error) {
	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	client := &Client{
		url:        defaultURL,
		httpClient: httpClient,
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) Handshake(req HandshakeRequest) error {
	const handshakePath = "/webhook/terraform"

	_, err := c.post(handshakePath, req)
	return err
}

func (c *Client) post(path string, req any) ([]byte, error) {
	var reqBody io.Reader

	if req != nil {
		payload, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		reqBody = bytes.NewBuffer(payload)
	}

	fullURL, err := url.JoinPath(c.url, path)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Post(fullURL, "application/json", reqBody)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error(
				fmt.Sprintf("error closing response body from %s", fullURL),
				"error", err.Error(),
			)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body from %s: %w", fullURL, err)
	}

	if resp.StatusCode/100 != 2 {
		return nil, &HTTPError{
			URL:        fullURL,
			StatusCode: resp.StatusCode,
			Body:       respBody,
		}
	}

	return respBody, nil
}
