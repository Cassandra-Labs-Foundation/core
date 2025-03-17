package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log" // Add this import
	"net/http"
	"time"
)

// Client provides methods to interact with the Supabase REST API
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Supabase client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// request makes an HTTP request to the Supabase API
func (c *Client) request(ctx context.Context, method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
    url := fmt.Sprintf("%s%s", c.baseURL, path)
    log.Printf("Making Supabase request: %s %s", method, url)

    var reqBody io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            log.Printf("Error marshaling request body: %v", err)
            return nil, fmt.Errorf("error marshaling request body: %w", err)
        }
        reqBody = bytes.NewBuffer(jsonBody)
        log.Printf("Request body: %s", string(jsonBody))
    }

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	// Add query parameters
	if queryParams != nil {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check if the response is successful
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("supabase API error: %s, status code: %d", string(respBody), resp.StatusCode)
	}

	return respBody, nil
}

// Insert inserts a record into the specified table
func (c *Client) Insert(ctx context.Context, table string, data interface{}) ([]byte, error) {
	return c.request(ctx, http.MethodPost, "/rest/v1/"+table, nil, data)
}

// Select retrieves records from the specified table
func (c *Client) Select(ctx context.Context, table string, queryParams map[string]string) ([]byte, error) {
	return c.request(ctx, http.MethodGet, "/rest/v1/"+table, queryParams, nil)
}

// SelectById retrieves a record from the specified table by its ID
func (c *Client) SelectById(ctx context.Context, table, id string) ([]byte, error) {
	queryParams := map[string]string{
		"id": fmt.Sprintf("eq.%s", id),
	}
	return c.request(ctx, http.MethodGet, "/rest/v1/"+table, queryParams, nil)
}

// Update updates a record in the specified table
func (c *Client) Update(ctx context.Context, table, id string, data interface{}) ([]byte, error) {
	queryParams := map[string]string{
		"id": fmt.Sprintf("eq.%s", id),
	}
	return c.request(ctx, http.MethodPatch, "/rest/v1/"+table, queryParams, data)
}

// Delete deletes a record from the specified table
func (c *Client) Delete(ctx context.Context, table, id string) ([]byte, error) {
	queryParams := map[string]string{
		"id": fmt.Sprintf("eq.%s", id),
	}
	return c.request(ctx, http.MethodDelete, "/rest/v1/"+table, queryParams, nil)
}