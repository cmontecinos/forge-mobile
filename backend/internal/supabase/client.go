// Package supabase provides a client wrapper for Supabase operations.
package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Supabase client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// AuthResponse represents Supabase auth response.
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	User         User   `json:"user"`
}

// User represents a Supabase user.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthError represents a Supabase auth error.
type AuthError struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

// NewClient creates a new Supabase client.
func NewClient(url, key string) *Client {
	return &Client{
		baseURL: url,
		apiKey:  key,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SignUp registers a new user with email and password.
func (c *Client) SignUp(email, password string) (*AuthResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	return c.authRequest("/auth/v1/signup", payload)
}

// SignIn authenticates a user with email and password.
func (c *Client) SignIn(email, password string) (*AuthResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	return c.authRequest("/auth/v1/token?grant_type=password", payload)
}

// RefreshToken refreshes an access token using a refresh token.
func (c *Client) RefreshToken(refreshToken string) (*AuthResponse, error) {
	payload := map[string]string{
		"refresh_token": refreshToken,
	}

	return c.authRequest("/auth/v1/token?grant_type=refresh_token", payload)
}

// SignOut invalidates a user's session.
func (c *Client) SignOut(accessToken string) error {
	req, err := http.NewRequest("POST", c.baseURL+"/auth/v1/logout", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("apikey", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("logout failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) authRequest(endpoint string, payload map[string]string) (*AuthResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var authErr AuthError
		if err := json.Unmarshal(respBody, &authErr); err != nil {
			return nil, fmt.Errorf("auth error: %s", string(respBody))
		}
		return nil, fmt.Errorf("%s: %s", authErr.Error, authErr.Description)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}
