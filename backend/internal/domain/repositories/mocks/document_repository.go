package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentRepository struct {
	mock.Mock
}

func (m *DocumentRepository) Create(ctx context.Context, document *entities.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *DocumentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Document), args.Error(1)
}

func (m *DocumentRepository) Update(ctx context.Context, document *entities.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *DocumentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *DocumentRepository) List(ctx context.Context, filter *repositories.DocumentFilter) ([]*entities.Document, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.Document), args.Get(1).(int64), args.Error(2)
}

func (m *DocumentRepository) GetByRelatedEntity(ctx context.Context, entityType string, entityID primitive.ObjectID) ([]*entities.Document, error) {
	args := m.Called(ctx, entityType, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Document), args.Error(1)
}

func (m *DocumentRepository) GetByType(ctx context.Context, docType entities.DocumentType) ([]*entities.Document, error) {
	args := m.Called(ctx, docType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Document), args.Error(1)
}

func (m *DocumentRepository) GetByUploader(ctx context.Context, userID primitive.ObjectID) ([]*entities.Document, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Document), args.Error(1)
}

func (m *DocumentRepository) GetPublicDocuments(ctx context.Context) ([]*entities.Document, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Document), args.Error(1)
}

func (m *DocumentRepository) GetExpiredDocuments(ctx context.Context) ([]*entities.Document, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Document), args.Error(1)
}

func (m *DocumentRepository) GetExpiringSoonDocuments(ctx context.Context) ([]*entities.Document, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Document), args.Error(1)
}

func (m *DocumentRepository) GetDocumentVersions(ctx context.Context, documentID primitive.ObjectID) ([]*entities.Document, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Document), args.Error(1)
}

func (m *DocumentRepository) IncrementDownloadCount(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *DocumentRepository) GetDocumentStatistics(ctx context.Context) (*repositories.DocumentStatistics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.DocumentStatistics), args.Error(1)
}

func (m *DocumentRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
