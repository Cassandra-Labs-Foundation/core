package business

import (
	"context"
	"errors"
	"time"

	"github.com/Cassandra-Labs-Foundation/core/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrInvalidBusiness  = errors.New("invalid business data")
	ErrBusinessNotFound = errors.New("business not found")
)

type Service interface {
	Create(ctx context.Context, input CreateBusinessInput) (*BusinessOutput, error)
	GetByID(ctx context.Context, id uuid.UUID) (*BusinessOutput, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateBusinessInput) (*BusinessOutput, error)
	List(ctx context.Context, limit, offset int) ([]*BusinessOutput, error)
}

// CreateBusinessInput represents the input for creating a business
type CreateBusinessInput struct {
	Name               string `json:"name" binding:"required"`
	RegistrationNumber string `json:"registration_number" binding:"required"`
	Address            string `json:"address" binding:"required"`
	Country            string `json:"country" binding:"required"`
	// New optional KYC fields:
	TaxID              *string `json:"tax_id"`
	KYCDocumentURL     *string `json:"kyc_document_url"`
}

// UpdateBusinessInput represents the input for updating a business
type UpdateBusinessInput struct {
	Name               *string `json:"name"`
	RegistrationNumber *string `json:"registration_number"`
	Address            *string `json:"address"`
	Country            *string `json:"country"`
	KYCStatus          *string `json:"kyc_status"`
	// New optional KYC fields:
	TaxID              *string `json:"tax_id"`
	KYCDocumentURL     *string `json:"kyc_document_url"`
}

type BusinessOutput struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	RegistrationNumber string     `json:"registration_number"`
	Address            string     `json:"address"`
	Country            string     `json:"country"`
	KYCStatus          string     `json:"kyc_status"`
	KYCVerifiedAt      *time.Time `json:"kyc_verified_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type service struct {
	businessRepo repository.BusinessRepository
}

func NewService(businessRepo repository.BusinessRepository) Service {
	return &service{
		businessRepo: businessRepo,
	}
}

func (s *service) Create(ctx context.Context, input CreateBusinessInput) (*BusinessOutput, error) {
	business := &repository.BusinessEntity{
		Name:               input.Name,
		RegistrationNumber: input.RegistrationNumber,
		Address:            input.Address,
		Country:            input.Country,
		KYCStatus:          "pending",
	}
	if err := s.businessRepo.Create(ctx, business); err != nil {
		return nil, err
	}
	return s.entityToOutput(business), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*BusinessOutput, error) {
	business, err := s.businessRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, ErrBusinessNotFound
	}
	return s.entityToOutput(business), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input UpdateBusinessInput) (*BusinessOutput, error) {
	business, err := s.businessRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, ErrBusinessNotFound
	}
	if input.Name != nil {
		business.Name = *input.Name
	}
	if input.RegistrationNumber != nil {
		business.RegistrationNumber = *input.RegistrationNumber
	}
	if input.Address != nil {
		business.Address = *input.Address
	}
	if input.Country != nil {
		business.Country = *input.Country
	}
	if input.KYCStatus != nil {
		business.KYCStatus = *input.KYCStatus
		if *input.KYCStatus == "verified" && business.KYCVerifiedAt == nil {
			now := time.Now()
			business.KYCVerifiedAt = &now
		}
	}
	if err := s.businessRepo.Update(ctx, business); err != nil {
		return nil, err
	}
	return s.entityToOutput(business), nil
}

func (s *service) List(ctx context.Context, limit, offset int) ([]*BusinessOutput, error) {
	businesses, err := s.businessRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	outputs := make([]*BusinessOutput, len(businesses))
	for i, b := range businesses {
		outputs[i] = s.entityToOutput(b)
	}
	return outputs, nil
}

func (s *service) entityToOutput(entity *repository.BusinessEntity) *BusinessOutput {
	return &BusinessOutput{
		ID:                 entity.ID,
		Name:               entity.Name,
		RegistrationNumber: entity.RegistrationNumber,
		Address:            entity.Address,
		Country:            entity.Country,
		KYCStatus:          entity.KYCStatus,
		KYCVerifiedAt:      entity.KYCVerifiedAt,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
	}
}
