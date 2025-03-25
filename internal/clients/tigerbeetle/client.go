package tigerbeetle

import (
	"context"
	"fmt"
)

// Client represents a TigerBeetle client.
type Client struct {
	Endpoint string
}

// NewClient creates a new TigerBeetle client instance.
func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

// CreateAccount simulates creating a new account in TigerBeetle.
func (c *Client) CreateAccount(ctx context.Context, accountID string, initialBalance int64) error {
	// In a real implementation, you would send an HTTP POST to c.Endpoint+"/account"
	// For now, we'll simulate it by printing the action and assuming success.
	fmt.Printf("Creating TigerBeetle account: ID=%s, Balance=%d\n", accountID, initialBalance)
	// Actual HTTP request code would go here.
	return nil
}

// Transfer simulates transferring funds between accounts.
func (c *Client) Transfer(ctx context.Context, fromAccountID, toAccountID string, amount int64) error {
	fmt.Printf("Transferring %d from %s to %s\n", amount, fromAccountID, toAccountID)
	// Actual HTTP request code would go here.
	return nil
}