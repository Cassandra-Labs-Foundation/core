package ledger

import (
	"context"

	"github.com/Cassandra-Labs-Foundation/core/internal/repository"
)

// Service defines ledger business operations.
type Service interface {
	CreateAccount(ctx context.Context, initialBalance int64) (string, error)
	TransferFunds(ctx context.Context, fromAccountID, toAccountID string, amount int64) error
}

type service struct {
	repo repository.LedgerRepository
}

// NewService creates a new ledger service.
func NewService(repo repository.LedgerRepository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateAccount(ctx context.Context, initialBalance int64) (string, error) {
	return s.repo.CreateAccount(ctx, initialBalance)
}

func (s *service) TransferFunds(ctx context.Context, fromAccountID, toAccountID string, amount int64) error {
	return s.repo.Transfer(ctx, fromAccountID, toAccountID, amount)
}