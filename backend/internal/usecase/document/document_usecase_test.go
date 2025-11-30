package document

import (
	"context"
	"testing"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDocumentUseCase_ArchiveDocument(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	useCase := NewDocumentUseCase(mockDocRepo, mockAuditLogRepo)

	docID := primitive.NewObjectID()
	uploaderID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	doc := &entities.Document{
		ID:         docID,
		UploadedBy: uploaderID,
	}

	t.Run("success", func(t *testing.T) {
		mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
		mockDocRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Document")).Run(func(args mock.Arguments) {
			arg := args.Get(1).(*entities.Document)
			assert.True(t, arg.IsArchived)
			assert.NotNil(t, arg.ArchivedAt)
		}).Return(nil).Once()
		mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil).Once()

		err := useCase.ArchiveDocument(context.Background(), docID, uploaderID)
		assert.NoError(t, err)
		mockDocRepo.AssertExpectations(t)
		mockAuditLogRepo.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()

		err := useCase.ArchiveDocument(context.Background(), docID, otherUserID)
		assert.Error(t, err)
		mockDocRepo.AssertExpectations(t)
	})
}

func TestDocumentUseCase_UnarchiveDocument(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	useCase := NewDocumentUseCase(mockDocRepo, mockAuditLogRepo)

	docID := primitive.NewObjectID()
	uploaderID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	now := time.Now()
	doc := &entities.Document{
		ID:         docID,
		UploadedBy: uploaderID,
		IsArchived: true,
		ArchivedAt: &now,
	}

	t.Run("success", func(t *testing.T) {
		mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()
		mockDocRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Document")).Run(func(args mock.Arguments) {
			arg := args.Get(1).(*entities.Document)
			assert.False(t, arg.IsArchived)
			assert.Nil(t, arg.ArchivedAt)
		}).Return(nil).Once()
		mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil).Once()

		err := useCase.UnarchiveDocument(context.Background(), docID, uploaderID)
		assert.NoError(t, err)
		mockDocRepo.AssertExpectations(t)
		mockAuditLogRepo.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil).Once()

		err := useCase.UnarchiveDocument(context.Background(), docID, otherUserID)
		assert.Error(t, err)
		mockDocRepo.AssertExpectations(t)
	})
}

func TestDocumentUseCase_SearchDocuments(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	useCase := NewDocumentUseCase(mockDocRepo, mockAuditLogRepo)

	userID := primitive.NewObjectID()
	docs := []*entities.Document{
		{ID: primitive.NewObjectID(), UploadedBy: userID},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID(), IsPublic: true},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID(), AccessibleBy: []primitive.ObjectID{userID}},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID()},
	}

	filter := &repositories.DocumentFilter{Search: "test"}

	mockDocRepo.On("List", mock.Anything, filter).Return(docs[:3], int64(3), nil).Once()

	accessibleDocs, total, err := useCase.SearchDocuments(context.Background(), filter, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, accessibleDocs, 3)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentUseCase_GetDocumentsByRelatedEntity(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	useCase := NewDocumentUseCase(mockDocRepo, nil)

	userID := primitive.NewObjectID()
	docs := []*entities.Document{
		{ID: primitive.NewObjectID(), UploadedBy: userID},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID(), IsPublic: true},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID(), AccessibleBy: []primitive.ObjectID{userID}},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID()},
	}

	mockDocRepo.On("GetByRelatedEntity", mock.Anything, "animal", mock.AnythingOfType("primitive.ObjectID")).Return(docs, nil).Once()

	accessibleDocs, err := useCase.GetDocumentsByRelatedEntity(context.Background(), "animal", primitive.NewObjectID(), userID)
	assert.NoError(t, err)
	assert.Len(t, accessibleDocs, 3)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentUseCase_GetDocumentsByType(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	useCase := NewDocumentUseCase(mockDocRepo, nil)

	userID := primitive.NewObjectID()
	docs := []*entities.Document{
		{ID: primitive.NewObjectID(), UploadedBy: userID},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID(), IsPublic: true},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID(), AccessibleBy: []primitive.ObjectID{userID}},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID()},
	}

	mockDocRepo.On("GetByType", mock.Anything, entities.DocumentType("medical")).Return(docs, nil).Once()

	accessibleDocs, err := useCase.GetDocumentsByType(context.Background(), entities.DocumentType("medical"), userID)
	assert.NoError(t, err)
	assert.Len(t, accessibleDocs, 3)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentUseCase_GetDocumentsByCategory(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	useCase := NewDocumentUseCase(mockDocRepo, nil)

	userID := primitive.NewObjectID()
	docs := []*entities.Document{
		{ID: primitive.NewObjectID(), UploadedBy: userID},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID(), IsPublic: true},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID(), AccessibleBy: []primitive.ObjectID{userID}},
		{ID: primitive.NewObjectID(), UploadedBy: primitive.NewObjectID()},
	}

	mockDocRepo.On("GetByType", mock.Anything, entities.DocumentType("legal")).Return(docs, nil).Once()

	accessibleDocs, err := useCase.GetDocumentsByCategory(context.Background(), entities.DocumentType("legal"), userID)
	assert.NoError(t, err)
	assert.Len(t, accessibleDocs, 3)
	mockDocRepo.AssertExpectations(t)
}
