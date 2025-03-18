package person

import (
	"context"
	"errors"
	"time"

	"github.com/Cassandra-Labs-Foundation/core/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrInvalidPerson = errors.New("invalid person data")
	ErrPersonNotFound = errors.New("person not found")
)

// Service provides person entity business logic
type Service interface {
	Create(ctx context.Context, input CreatePersonInput) (*PersonOutput, error)
	GetByID(ctx context.Context, id uuid.UUID) (*PersonOutput, error)
	Update(ctx context.Context, id uuid.UUID, input UpdatePersonInput) (*PersonOutput, error)
	List(ctx context.Context, limit, offset int) ([]*PersonOutput, error)
}

// CreatePersonInput represents the input for creating a person
type CreatePersonInput struct {
	FirstName      string  `json:"first_name" binding:"required"`
	LastName       string  `json:"last_name" binding:"required"`
	DateOfBirth    string  `json:"date_of_birth" binding:"required"` // Format: YYYY-MM-DD
	SSN            *string `json:"ssn"`
	Email          *string `json:"email"`
	PhoneNumber    *string `json:"phone_number"`
	Street1        *string `json:"street1"`
	Street2        *string `json:"street2"`
	City           *string `json:"city"`
	State          *string `json:"state"`
	PostalCode     *string `json:"postal_code"`
	Country        *string `json:"country"`
	// New optional KYC fields:
	GovernmentID   *string `json:"government_id"`
	Nationality    *string `json:"nationality"`
	KYCDocumentURL *string `json:"kyc_document_url"`
}

// UpdatePersonInput represents the input for updating a person
type UpdatePersonInput struct {
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	DateOfBirth    *string `json:"date_of_birth"` // Format: YYYY-MM-DD
	SSN            *string `json:"ssn"`
	Email          *string `json:"email"`
	PhoneNumber    *string `json:"phone_number"`
	Street1        *string `json:"street1"`
	Street2        *string `json:"street2"`
	City           *string `json:"city"`
	State          *string `json:"state"`
	PostalCode     *string `json:"postal_code"`
	Country        *string `json:"country"`
	KYCStatus      *string `json:"kyc_status"`
	// New optional KYC fields:
	GovernmentID   *string `json:"government_id"`
	Nationality    *string `json:"nationality"`
	KYCDocumentURL *string `json:"kyc_document_url"`
}

// PersonOutput represents the output for person entity operations
type PersonOutput struct {
	ID            uuid.UUID  `json:"id"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	DateOfBirth   string     `json:"date_of_birth"` // Format: YYYY-MM-DD
	SSN           *string    `json:"ssn,omitempty"`
	Email         *string    `json:"email,omitempty"`
	PhoneNumber   *string    `json:"phone_number,omitempty"`
	Street1       *string    `json:"street1,omitempty"`
	Street2       *string    `json:"street2,omitempty"`
	City          *string    `json:"city,omitempty"`
	State         *string    `json:"state,omitempty"`
	PostalCode    *string    `json:"postal_code,omitempty"`
	Country       *string    `json:"country,omitempty"`
	KYCStatus     string     `json:"kyc_status"`
	KYCVerifiedAt *time.Time `json:"kyc_verified_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type service struct {
	personRepo repository.PersonRepository
}

// NewService creates a new person service
func NewService(personRepo repository.PersonRepository) Service {
	return &service{
		personRepo: personRepo,
	}
}

// Create creates a new person entity
func (s *service) Create(ctx context.Context, input CreatePersonInput) (*PersonOutput, error) {
	// Parse date of birth
	dob, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		return nil, ErrInvalidPerson
	}

	// Create person entity with new KYC fields
	person := &repository.PersonEntity{
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		DateOfBirth:     dob,
		SSN:             input.SSN,
		Email:           input.Email,
		PhoneNumber:     input.PhoneNumber,
		Street1:         input.Street1,
		Street2:         input.Street2,
		City:            input.City,
		State:           input.State,
		PostalCode:      input.PostalCode,
		Country:         input.Country,
		KYCStatus:       "pending", // Default KYC status
		// New fields:
		GovernmentID:    input.GovernmentID,
		Nationality:     input.Nationality,
		KYCDocumentURL:  input.KYCDocumentURL,
	}

	if err := s.personRepo.Create(ctx, person); err != nil {
		return nil, err
	}

	return s.entityToOutput(person), nil
}

// GetByID retrieves a person entity by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*PersonOutput, error) {
	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if person == nil {
		return nil, ErrPersonNotFound
	}

	return s.entityToOutput(person), nil
}

// Update updates an existing person entity
func (s *service) Update(ctx context.Context, id uuid.UUID, input UpdatePersonInput) (*PersonOutput, error) {
	// Get existing person
	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if person == nil {
		return nil, ErrPersonNotFound
	}

	// Update fields if provided
	if input.FirstName != nil {
		person.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		person.LastName = *input.LastName
	}
	if input.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *input.DateOfBirth)
		if err != nil {
			return nil, ErrInvalidPerson
		}
		person.DateOfBirth = dob
	}
	if input.SSN != nil {
		person.SSN = input.SSN
	}
	if input.Email != nil {
		person.Email = input.Email
	}
	if input.PhoneNumber != nil {
		person.PhoneNumber = input.PhoneNumber
	}
	if input.Street1 != nil {
		person.Street1 = input.Street1
	}
	if input.Street2 != nil {
		person.Street2 = input.Street2
	}
	if input.City != nil {
		person.City = input.City
	}
	if input.State != nil {
		person.State = input.State
	}
	if input.PostalCode != nil {
		person.PostalCode = input.PostalCode
	}
	if input.Country != nil {
		person.Country = input.Country
	}
	if input.KYCStatus != nil {
		person.KYCStatus = *input.KYCStatus
		
		// If status changed to verified, update verification timestamp
		if *input.KYCStatus == "verified" && (person.KYCVerifiedAt == nil) {
			now := time.Now()
			person.KYCVerifiedAt = &now
		}
	}

	// Update in database
	if err := s.personRepo.Update(ctx, person); err != nil {
		return nil, err
	}

	return s.entityToOutput(person), nil
}

// List retrieves a paginated list of person entities
func (s *service) List(ctx context.Context, limit, offset int) ([]*PersonOutput, error) {
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

	people, err := s.personRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	outputs := make([]*PersonOutput, len(people))
	for i, person := range people {
		outputs[i] = s.entityToOutput(person)
	}

	return outputs, nil
}

// Helper function to convert entity to output
func (s *service) entityToOutput(entity *repository.PersonEntity) *PersonOutput {
	return &PersonOutput{
		ID:             entity.ID,
		FirstName:      entity.FirstName,
		LastName:       entity.LastName,
		DateOfBirth:    entity.DateOfBirth.Format("2006-01-02"),
		SSN:            entity.SSN,
		Email:          entity.Email,
		PhoneNumber:    entity.PhoneNumber,
		Street1:        entity.Street1,
		Street2:        entity.Street2,
		City:           entity.City,
		State:          entity.State,
		PostalCode:     entity.PostalCode,
		Country:        entity.Country,
		KYCStatus:      entity.KYCStatus,
		KYCVerifiedAt:  entity.KYCVerifiedAt,
		// New fields:
		// (Assuming you add these to PersonOutput as well)
		// GovernmentID:    entity.GovernmentID,
		// Nationality:     entity.Nationality,
		// KYCDocumentURL:  entity.KYCDocumentURL,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}