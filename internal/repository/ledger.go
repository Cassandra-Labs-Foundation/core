package repository

import (
	"context"
	"errors"

	"github.com/Cassandra-Labs-Foundation/core/internal/clients/tigerbeetle"
	"github.com/google/uuid"
)

// LedgerRepository defines methods for ledger operations.
type LedgerRepository interface {
	CreateAccount(ctx context.Context, initialBalance int64) (string, error)
	Transfer(ctx context.Context, fromAccountID, toAccountID string, amount int64) error
}

type ledgerRepository struct {
	client *tigerbeetle.Client
}

// NewLedgerRepository creates a new LedgerRepository.
func NewLedgerRepository(client *tigerbeetle.Client) LedgerRepository {
	return &ledgerRepository{
		client: client,
	}
}

// CreateAccount creates a new account in TigerBeetle.
// It generates a new UUID for the account.
func (r *ledgerRepository) CreateAccount(ctx context.Context, initialBalance int64) (string, error) {
	accountID := uuid.New().String()
	err := r.client.CreateAccount(ctx, accountID, initialBalance)
	if err != nil {
		return "", err
	}
	return accountID, nil
}

// Transfer executes a fund transfer between two accounts.
func (r *ledgerRepository) Transfer(ctx context.Context, fromAccountID, toAccountID string, amount int64) error {
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}
	return r.client.Transfer(ctx, fromAccountID, toAccountID, amount)
}