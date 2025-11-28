package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/usecase/document"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDocumentHandler_ArchiveDocument(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	docUseCase := document.NewDocumentUseCase(mockDocRepo, mockAuditLogRepo)
	handler := NewDocumentHandler(docUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	uploaderID := primitive.NewObjectID()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uploaderID)
	})
	router.POST("/documents/:id/archive", handler.ArchiveDocument)

	docID := primitive.NewObjectID()
	doc := &entities.Document{ID: docID, UploadedBy: uploaderID}

	mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil)
	mockDocRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockAuditLogRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/documents/"+docID.Hex()+"/archive", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDocumentHandler_UnarchiveDocument(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	docUseCase := document.NewDocumentUseCase(mockDocRepo, mockAuditLogRepo)
	handler := NewDocumentHandler(docUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	uploaderID := primitive.NewObjectID()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uploaderID)
	})
	router.POST("/documents/:id/unarchive", handler.UnarchiveDocument)

	docID := primitive.NewObjectID()
	doc := &entities.Document{ID: docID, UploadedBy: uploaderID}

	mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil)
	mockDocRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockAuditLogRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/documents/"+docID.Hex()+"/unarchive", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDocumentHandler_ShareDocument(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	docUseCase := document.NewDocumentUseCase(mockDocRepo, mockAuditLogRepo)
	handler := NewDocumentHandler(docUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	uploaderID := primitive.NewObjectID()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uploaderID)
	})
	router.POST("/documents/:id/share", handler.ShareDocument)

	docID := primitive.NewObjectID()
	targetUserID := primitive.NewObjectID()
	doc := &entities.Document{ID: docID, UploadedBy: uploaderID}

	mockDocRepo.On("FindByID", mock.Anything, docID).Return(doc, nil)
	mockDocRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockAuditLogRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(gin.H{"user_id": targetUserID.Hex()})
	req, _ := http.NewRequest("POST", "/documents/"+docID.Hex()+"/share", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDocumentHandler_SearchDocuments(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDocRepo := new(mocks.DocumentRepository)
		docUseCase := document.NewDocumentUseCase(mockDocRepo, nil)
		handler := NewDocumentHandler(docUseCase)

		router := setupRouter()
		router.GET("/documents/search", handler.SearchDocuments)

		docs := []*entities.Document{{ID: primitive.NewObjectID()}}
		mockDocRepo.On("List", mock.Anything, mock.AnythingOfType("*repositories.DocumentFilter")).Return(docs, int64(1), nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/documents/search?q=test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp gin.H
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(1), resp["total"])
	})

	t.Run("missing q", func(t *testing.T) {
		mockDocRepo := new(mocks.DocumentRepository)
		docUseCase := document.NewDocumentUseCase(mockDocRepo, nil)
		handler := NewDocumentHandler(docUseCase)

		router := setupRouter()
		router.GET("/documents/search", handler.SearchDocuments)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/documents/search", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDocumentHandler_GetDocumentsByEntity(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	docUseCase := document.NewDocumentUseCase(mockDocRepo, nil)
	handler := NewDocumentHandler(docUseCase)

	router := setupRouter()
	router.GET("/documents/entity/:entity_type/:entity_id", handler.GetDocumentsByEntity)

	docs := []*entities.Document{{ID: primitive.NewObjectID()}}
	mockDocRepo.On("GetByRelatedEntity", mock.Anything, "animal", mock.AnythingOfType("primitive.ObjectID")).Return(docs, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/documents/entity/animal/"+primitive.NewObjectID().Hex(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDocumentHandler_GetDocumentsByType(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	docUseCase := document.NewDocumentUseCase(mockDocRepo, nil)
	handler := NewDocumentHandler(docUseCase)

	router := setupRouter()
	router.GET("/documents/type/:type", handler.GetDocumentsByType)

	docs := []*entities.Document{{ID: primitive.NewObjectID()}}
	mockDocRepo.On("GetByType", mock.Anything, entities.DocumentType("medical")).Return(docs, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/documents/type/medical", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDocumentHandler_GetDocumentsByCategory(t *testing.T) {
	mockDocRepo := new(mocks.DocumentRepository)
	docUseCase := document.NewDocumentUseCase(mockDocRepo, nil)
	handler := NewDocumentHandler(docUseCase)

	router := setupRouter()
	router.GET("/documents/category/:category", handler.GetDocumentsByCategory)

	docs := []*entities.Document{{ID: primitive.NewObjectID()}}
	mockDocRepo.On("GetByType", mock.Anything, entities.DocumentType("legal")).Return(docs, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/documents/category/legal", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
