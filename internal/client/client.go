// Copyright sugarshin 2026
// SPDX-License-Identifier: MIT

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/line/line-bot-sdk-go/v8/linebot/liff"
)

type Client struct {
	liffAPI *liff.LiffAPI

	// For stateless token management
	channelID     string
	channelSecret string
	tokenMu       sync.Mutex
	token         string
	tokenExpiry   time.Time

	// Options for recreating liffAPI with fresh token
	options []liff.LiffAPIOption
}

func NewWithToken(token string, options ...liff.LiffAPIOption) (*Client, error) {
	api, err := liff.NewLiffAPI(token, options...)
	if err != nil {
		return nil, fmt.Errorf("creating LIFF API client: %w", err)
	}
	return &Client{
		liffAPI: api,
		token:   token,
		options: options,
	}, nil
}

func NewWithCredentials(channelID, channelSecret string, options ...liff.LiffAPIOption) (*Client, error) {
	c := &Client{
		channelID:     channelID,
		channelSecret: channelSecret,
		options:       options,
	}
	if err := c.refreshToken(); err != nil {
		return nil, fmt.Errorf("issuing stateless channel access token: %w", err)
	}
	return c, nil
}

func (c *Client) ensureToken() error {
	if c.channelID == "" {
		return nil
	}
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	if time.Now().Before(c.tokenExpiry.Add(-1 * time.Minute)) {
		return nil
	}
	return c.refreshToken()
}

func (c *Client) refreshToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.channelID)
	data.Set("client_secret", c.channelSecret)

	resp, err := http.PostForm("https://api.line.me/oauth2/v3/token", data)
	if err != nil {
		return fmt.Errorf("requesting stateless token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("stateless token request failed with status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decoding token response: %w", err)
	}

	c.token = result.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	api, err := liff.NewLiffAPI(c.token, c.options...)
	if err != nil {
		return fmt.Errorf("creating LIFF API client with new token: %w", err)
	}
	c.liffAPI = api
	return nil
}

func (c *Client) AddLiffApp(ctx context.Context, req *liff.AddLiffAppRequest) (*liff.AddLiffAppResponse, error) {
	if err := c.ensureToken(); err != nil {
		return nil, err
	}
	return c.liffAPI.WithContext(ctx).AddLIFFApp(req)
}

func (c *Client) GetAllLiffApps(ctx context.Context) ([]liff.LiffApp, error) {
	if err := c.ensureToken(); err != nil {
		return nil, err
	}
	resp, err := c.liffAPI.WithContext(ctx).GetAllLIFFApps()
	if err != nil {
		return nil, err
	}
	return resp.Apps, nil
}

func (c *Client) GetLiffApp(ctx context.Context, liffID string) (*liff.LiffApp, error) {
	apps, err := c.GetAllLiffApps(ctx)
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		if app.LiffId == liffID {
			return &app, nil
		}
	}
	return nil, fmt.Errorf("LIFF app %s not found", liffID)
}

func (c *Client) UpdateLiffApp(ctx context.Context, liffID string, req *liff.UpdateLiffAppRequest) error {
	if err := c.ensureToken(); err != nil {
		return err
	}
	_, err := c.liffAPI.WithContext(ctx).UpdateLIFFApp(liffID, req)
	return err
}

func (c *Client) DeleteLiffApp(ctx context.Context, liffID string) error {
	if err := c.ensureToken(); err != nil {
		return err
	}
	_, err := c.liffAPI.WithContext(ctx).DeleteLIFFApp(liffID)
	return err
}

// IsNotFound checks if the error indicates a LIFF app was not found.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "not found")
}
