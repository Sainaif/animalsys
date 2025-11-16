package contact

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UseCase provides contact business logic
type UseCase struct {
	repo repositories.ContactRepository
}

// NewUseCase creates contact usecase
func NewUseCase(repo repositories.ContactRepository) *UseCase {
	return &UseCase{repo: repo}
}

// List retrieves contacts by filter
func (uc *UseCase) List(ctx context.Context, filter repositories.ContactFilter) ([]*entities.Contact, int64, error) {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return uc.repo.List(ctx, filter)
}

// Get contact by id
func (uc *UseCase) Get(ctx context.Context, id primitive.ObjectID) (*entities.Contact, error) {
	return uc.repo.FindByID(ctx, id)
}

// Create request payload
type CreateRequest struct {
	FirstName        string                 `json:"first_name" validate:"required"`
	LastName         string                 `json:"last_name" validate:"required"`
	Organization     string                 `json:"organization,omitempty"`
	Email            string                 `json:"email,omitempty" validate:"omitempty,email"`
	Phone            string                 `json:"phone,omitempty"`
	Type             entities.ContactType   `json:"type" validate:"required"`
	Status           entities.ContactStatus `json:"status" validate:"required"`
	Tags             []string               `json:"tags,omitempty"`
	OwnerID          string                 `json:"owner_id,omitempty"`
	OwnerName        string                 `json:"owner_name,omitempty"`
	PreferredChannel string                 `json:"preferred_channel,omitempty"`
	NextFollowUpAt   *time.Time             `json:"next_follow_up_at,omitempty"`
	Notes            string                 `json:"notes,omitempty"`
	Address          *entities.AddressInfo  `json:"address,omitempty"`
}

// UpdateRequest for patch
type UpdateRequest struct {
	FirstName        *string                 `json:"first_name,omitempty"`
	LastName         *string                 `json:"last_name,omitempty"`
	Organization     *string                 `json:"organization,omitempty"`
	Email            *string                 `json:"email,omitempty"`
	Phone            *string                 `json:"phone,omitempty"`
	Type             *entities.ContactType   `json:"type,omitempty"`
	Status           *entities.ContactStatus `json:"status,omitempty"`
	Tags             *[]string               `json:"tags,omitempty"`
	OwnerID          *string                 `json:"owner_id,omitempty"`
	OwnerName        *string                 `json:"owner_name,omitempty"`
	PreferredChannel *string                 `json:"preferred_channel,omitempty"`
	NextFollowUpAt   **time.Time             `json:"next_follow_up_at,omitempty"`
	Notes            *string                 `json:"notes,omitempty"`
	Address          *entities.AddressInfo   `json:"address,omitempty"`
}

// Create contact
func (uc *UseCase) Create(ctx context.Context, req *CreateRequest) (*entities.Contact, error) {
	if req.FirstName == "" || req.LastName == "" {
		return nil, errors.NewBadRequest("first and last name are required")
	}

	contact := &entities.Contact{
		ID:               primitive.NewObjectID(),
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Organization:     req.Organization,
		Email:            req.Email,
		Phone:            req.Phone,
		Type:             req.Type,
		Status:           req.Status,
		Tags:             req.Tags,
		PreferredChannel: req.PreferredChannel,
		Notes:            req.Notes,
		Address:          req.Address,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if req.OwnerID != "" {
		ownerID, err := primitive.ObjectIDFromHex(req.OwnerID)
		if err != nil {
			return nil, errors.NewBadRequest("invalid owner id")
		}
		contact.OwnerID = &ownerID
	}
	if req.OwnerName != "" {
		contact.OwnerName = req.OwnerName
	}

	if req.NextFollowUpAt != nil {
		contact.NextFollowUpAt = req.NextFollowUpAt
	}

	if err := uc.repo.Create(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// Update contact
func (uc *UseCase) Update(ctx context.Context, id primitive.ObjectID, req *UpdateRequest) (*entities.Contact, error) {
	contact, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		contact.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		contact.LastName = *req.LastName
	}
	if req.Organization != nil {
		contact.Organization = *req.Organization
	}
	if req.Email != nil {
		contact.Email = *req.Email
	}
	if req.Phone != nil {
		contact.Phone = *req.Phone
	}
	if req.Type != nil {
		contact.Type = *req.Type
	}
	if req.Status != nil {
		contact.Status = *req.Status
	}
	if req.Tags != nil {
		contact.Tags = *req.Tags
	}
	if req.OwnerID != nil {
		if *req.OwnerID == "" {
			contact.OwnerID = nil
		} else {
			ownerID, convErr := primitive.ObjectIDFromHex(*req.OwnerID)
			if convErr != nil {
				return nil, errors.NewBadRequest("invalid owner id")
			}
			contact.OwnerID = &ownerID
		}
	}
	if req.OwnerName != nil {
		contact.OwnerName = *req.OwnerName
	}
	if req.PreferredChannel != nil {
		contact.PreferredChannel = *req.PreferredChannel
	}
	if req.Notes != nil {
		contact.Notes = *req.Notes
	}
	if req.Address != nil {
		contact.Address = req.Address
	}
	if req.NextFollowUpAt != nil {
		contact.NextFollowUpAt = *req.NextFollowUpAt
	}

	contact.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// Delete removes contact
func (uc *UseCase) Delete(ctx context.Context, id primitive.ObjectID) error {
	return uc.repo.Delete(ctx, id)
}
