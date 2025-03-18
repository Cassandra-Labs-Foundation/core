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

// PersonEntity represents a person entity in the database
// PersonEntity represents a person entity in the database
type PersonEntity struct {
	ID              uuid.UUID  `json:"id,omitempty"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	DateOfBirth     time.Time  `json:"date_of_birth"`
	SSN             *string    `json:"ssn,omitempty"`
	Email           *string    `json:"email,omitempty"`
	PhoneNumber     *string    `json:"phone_number,omitempty"`
	Street1         *string    `json:"street1,omitempty"`
	Street2         *string    `json:"street2,omitempty"`
	City            *string    `json:"city,omitempty"`
	State           *string    `json:"state,omitempty"`
	PostalCode      *string    `json:"postal_code,omitempty"`
	Country         *string    `json:"country,omitempty"`
	KYCStatus       string     `json:"kyc_status,omitempty"`
	KYCVerifiedAt   *time.Time `json:"kyc_verified_at,omitempty"`
	// New KYC fields:
	GovernmentID    *string    `json:"government_id,omitempty"`
	Nationality     *string    `json:"nationality,omitempty"`
	KYCDocumentURL  *string    `json:"kyc_document_url,omitempty"`
	CreatedAt       time.Time  `json:"created_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at,omitempty"`
}

// PersonRepository provides methods to interact with person entities in the database
type PersonRepository interface {
	Create(ctx context.Context, person *PersonEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*PersonEntity, error)
	Update(ctx context.Context, person *PersonEntity) error
	List(ctx context.Context, limit, offset int) ([]*PersonEntity, error)
}

type personRestRepository struct {
	client *supabase.Client
	table  string
}

// NewPersonRestRepository creates a new person repository using Supabase REST API
func NewPersonRestRepository(client *supabase.Client) PersonRepository {
	return &personRestRepository{
		client: client,
		table:  "person_entities",
	}
}

// Create inserts a new person entity
func (r *personRestRepository) Create(ctx context.Context, person *PersonEntity) error {
	// Set default KYC status if not provided
	if person.KYCStatus == "" {
		person.KYCStatus = "pending"
	}

	// Build the payload map.
	// Format date_of_birth as "YYYY-MM-DD" to match the database type (date)
	payload := map[string]interface{}{
		"first_name":     person.FirstName,
		"last_name":      person.LastName,
		"date_of_birth":  person.DateOfBirth.Format("2006-01-02"),
		"ssn":            person.SSN,
		"email":          person.Email,
		"phone_number":   person.PhoneNumber,
		"street1":        person.Street1,
		"street2":        person.Street2,
		"city":           person.City,
		"state":          person.State,
		"postal_code":    person.PostalCode,
		"country":        person.Country,
		"kyc_status":     person.KYCStatus,
		"kyc_verified_at": person.KYCVerifiedAt,
	}

	// Only include the id if it is not zero.
	if person.ID != uuid.Nil {
		payload["id"] = person.ID
	}

	// Insert the person entity using the payload map
	respBody, err := r.client.Insert(ctx, r.table, payload)
	if err != nil {
		return err
	}

	// Parse the response to update the entity with generated fields
	var createdPersons []*PersonEntity
	if err := json.Unmarshal(respBody, &createdPersons); err != nil {
		return err
	}

	if len(createdPersons) == 0 {
		return errors.New("no person entity was created")
	}

	// Update the input entity with the created entity's data
	createdPerson := createdPersons[0]
	person.ID = createdPerson.ID
	person.CreatedAt = createdPerson.CreatedAt
	person.UpdatedAt = createdPerson.UpdatedAt

	return nil
}

// GetByID retrieves a person entity by its ID
func (r *personRestRepository) GetByID(ctx context.Context, id uuid.UUID) (*PersonEntity, error) {
	respBody, err := r.client.SelectById(ctx, r.table, id.String())
	if err != nil {
		return nil, err
	}

	var persons []*PersonEntity
	if err := json.Unmarshal(respBody, &persons); err != nil {
		return nil, err
	}

	if len(persons) == 0 {
		return nil, nil // Not found
	}

	return persons[0], nil
}

// Update updates an existing person entity
func (r *personRestRepository) Update(ctx context.Context, person *PersonEntity) error {
	if person.ID == uuid.Nil {
		return errors.New("person ID is required for update")
	}

	// Build the update payload.
	updateData := map[string]interface{}{
		"first_name":      person.FirstName,
		"last_name":       person.LastName,
		"date_of_birth":   person.DateOfBirth,
		"ssn":             person.SSN,
		"email":           person.Email,
		"phone_number":    person.PhoneNumber,
		"street1":         person.Street1,
		"street2":         person.Street2,
		"city":            person.City,
		"state":           person.State,
		"postal_code":     person.PostalCode,
		"country":         person.Country,
		"kyc_status":      person.KYCStatus,
		"kyc_verified_at": person.KYCVerifiedAt,
	}

	// Update the person entity
	respBody, err := r.client.Update(ctx, r.table, person.ID.String(), updateData)
	if err != nil {
		return err
	}

	// Parse the response to update the entity with updated fields
	var updatedPersons []*PersonEntity
	if err := json.Unmarshal(respBody, &updatedPersons); err != nil {
		return err
	}

	if len(updatedPersons) == 0 {
		return errors.New("no person entity was updated")
	}

	// Update the input entity with the updated entity's data
	updatedPerson := updatedPersons[0]
	person.UpdatedAt = updatedPerson.UpdatedAt

	return nil
}

// List retrieves a paginated list of person entities
func (r *personRestRepository) List(ctx context.Context, limit, offset int) ([]*PersonEntity, error) {
	// Apply sensible defaults
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Prepare query parameters
	queryParams := map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(offset),
		"order":  "created_at.desc",
	}

	// Fetch the person entities
	respBody, err := r.client.Select(ctx, r.table, queryParams)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var persons []*PersonEntity
	if err := json.Unmarshal(respBody, &persons); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return persons, nil
}