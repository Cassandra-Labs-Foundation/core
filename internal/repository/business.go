package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Cassandra-Labs-Foundation/core/internal/clients/supabase"
	"github.com/google/uuid"
)

// BusinessEntity represents a business entity in the database
type BusinessEntity struct {
	ID                 uuid.UUID  `json:"id,omitempty"`
	Name               string     `json:"name"`
	RegistrationNumber string     `json:"registration_number"`
	Address            string     `json:"address"`
	Country            string     `json:"country"`
	KYCStatus          string     `json:"kyc_status,omitempty"`
	KYCVerifiedAt      *time.Time `json:"kyc_verified_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at,omitempty"`
	UpdatedAt          time.Time  `json:"updated_at,omitempty"`
}

// BusinessRepository provides methods to interact with business entities
type BusinessRepository interface {
	Create(ctx context.Context, business *BusinessEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*BusinessEntity, error)
	Update(ctx context.Context, business *BusinessEntity) error
	List(ctx context.Context, limit, offset int) ([]*BusinessEntity, error)
}

type businessRestRepository struct {
	client *supabase.Client
	table  string
}

// NewBusinessRestRepository creates a new business repository using Supabase REST API
func NewBusinessRestRepository(client *supabase.Client) BusinessRepository {
	return &businessRestRepository{
		client: client,
		table:  "business_entities",
	}
}

func (r *businessRestRepository) Create(ctx context.Context, business *BusinessEntity) error {
	if business.KYCStatus == "" {
		business.KYCStatus = "pending"
	}

	payload := map[string]interface{}{
		"name":                business.Name,
		"registration_number": business.RegistrationNumber,
		"address":             business.Address,
		"country":             business.Country,
		"kyc_status":          business.KYCStatus,
		"kyc_verified_at":     business.KYCVerifiedAt,
	}

	// Only include the ID if already set (non-zero)
	if business.ID != uuid.Nil {
		payload["id"] = business.ID
	}

	respBody, err := r.client.Insert(ctx, r.table, payload)
	if err != nil {
		return err
	}

	var createdBusinesses []*BusinessEntity
	if err := json.Unmarshal(respBody, &createdBusinesses); err != nil {
		return err
	}
	if len(createdBusinesses) == 0 {
		return errors.New("no business entity was created")
	}

	created := createdBusinesses[0]
	business.ID = created.ID
	business.CreatedAt = created.CreatedAt
	business.UpdatedAt = created.UpdatedAt

	return nil
}

func (r *businessRestRepository) GetByID(ctx context.Context, id uuid.UUID) (*BusinessEntity, error) {
	respBody, err := r.client.SelectById(ctx, r.table, id.String())
	if err != nil {
		return nil, err
	}
	var businesses []*BusinessEntity
	if err := json.Unmarshal(respBody, &businesses); err != nil {
		return nil, err
	}
	if len(businesses) == 0 {
		return nil, nil
	}
	return businesses[0], nil
}

func (r *businessRestRepository) Update(ctx context.Context, business *BusinessEntity) error {
	if business.ID == uuid.Nil {
		return errors.New("business ID is required for update")
	}

	payload := map[string]interface{}{
		"name":                business.Name,
		"registration_number": business.RegistrationNumber,
		"address":             business.Address,
		"country":             business.Country,
		"kyc_status":          business.KYCStatus,
		"kyc_verified_at":     business.KYCVerifiedAt,
	}

	respBody, err := r.client.Update(ctx, r.table, business.ID.String(), payload)
	if err != nil {
		return err
	}

	var updatedBusinesses []*BusinessEntity
	if err := json.Unmarshal(respBody, &updatedBusinesses); err != nil {
		return err
	}
	if len(updatedBusinesses) == 0 {
		return errors.New("no business entity was updated")
	}

	updated := updatedBusinesses[0]
	business.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *businessRestRepository) List(ctx context.Context, limit, offset int) ([]*BusinessEntity, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	queryParams := map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(offset),
		"order":  "created_at.desc",
	}
	respBody, err := r.client.Select(ctx, r.table, queryParams)
	if err != nil {
		return nil, err
	}
	var businesses []*BusinessEntity
	if err := json.Unmarshal(respBody, &businesses); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}
	return businesses, nil
}