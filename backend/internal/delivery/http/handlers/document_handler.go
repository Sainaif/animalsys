package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/document"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentHandler handles document-related HTTP requests
type DocumentHandler struct {
	documentUseCase *document.DocumentUseCase
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler(documentUseCase *document.DocumentUseCase) *DocumentHandler {
	return &DocumentHandler{
		documentUseCase: documentUseCase,
	}
}

// CreateDocument creates a new document
func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var docReq entities.Document
	if err := c.ShouldBindJSON(&docReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.CreateDocument(c.Request.Context(), &docReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, docReq)
}

// GetDocument gets a document by ID
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	doc, err := h.documentUseCase.GetDocumentByID(c.Request.Context(), id, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, doc)
}

// UpdateDocument updates a document
func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var docReq entities.Document
	if err := c.ShouldBindJSON(&docReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docReq.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.UpdateDocument(c.Request.Context(), &docReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, docReq)
}

// DeleteDocument deletes a document
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.DeleteDocument(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// ListDocuments lists documents with filtering
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	filter := &repositories.DocumentFilter{}

	// Parse query parameters
	filter.Type = c.Query("type")
	filter.RelatedEntity = c.Query("related_entity")
	filter.Search = c.Query("search")

	if relatedIDStr := c.Query("related_entity_id"); relatedIDStr != "" {
		relatedID, err := primitive.ObjectIDFromHex(relatedIDStr)
		if err == nil {
			filter.RelatedEntityID = &relatedID
		}
	}

	if uploadedByStr := c.Query("uploaded_by"); uploadedByStr != "" {
		uploadedBy, err := primitive.ObjectIDFromHex(uploadedByStr)
		if err == nil {
			filter.UploadedBy = &uploadedBy
		}
	}

	if publicStr := c.Query("is_public"); publicStr != "" {
		isPublic := publicStr == "true"
		filter.IsPublic = &isPublic
	}

	if expiredStr := c.Query("is_expired"); expiredStr != "" {
		isExpired := expiredStr == "true"
		filter.IsExpired = &isExpired
	}

	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filter.Tags = tags
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "created_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	userID := c.MustGet("user_id").(primitive.ObjectID)

	documents, total, err := h.documentUseCase.SearchDocuments(c.Request.Context(), filter, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   documents,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetPublicDocuments gets all public documents
func (h *DocumentHandler) GetPublicDocuments(c *gin.Context) {
	documents, err := h.documentUseCase.GetPublicDocuments(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, documents)
}

// GetMyDocuments gets documents uploaded by the current user
func (h *DocumentHandler) GetMyDocuments(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	documents, err := h.documentUseCase.GetDocumentsByUploader(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, documents)
}

// GetExpiredDocuments gets all expired documents
func (h *DocumentHandler) GetExpiredDocuments(c *gin.Context) {
	documents, err := h.documentUseCase.GetExpiredDocuments(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, documents)
}

// GetExpiringSoonDocuments gets documents expiring soon
func (h *DocumentHandler) GetExpiringSoonDocuments(c *gin.Context) {
	documents, err := h.documentUseCase.GetExpiringSoonDocuments(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, documents)
}

// GetDocumentVersions gets all versions of a document
func (h *DocumentHandler) GetDocumentVersions(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	versions, err := h.documentUseCase.GetDocumentVersions(c.Request.Context(), id, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, versions)
}

// GetDocumentStatistics gets document statistics
func (h *DocumentHandler) GetDocumentStatistics(c *gin.Context) {
	stats, err := h.documentUseCase.GetDocumentStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DownloadDocument handles document download
func (h *DocumentHandler) DownloadDocument(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	doc, err := h.documentUseCase.DownloadDocument(c.Request.Context(), id, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	// In a real application, you would serve the actual file here
	// For now, we return the document metadata
	c.JSON(http.StatusOK, gin.H{
		"document":     doc,
		"download_url": doc.FileURL, // Use the document's file URL
	})
}

// CreateNewVersion creates a new version of a document
func (h *DocumentHandler) CreateNewVersion(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req struct {
		FilePath string `json:"file_path" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	newVersion, err := h.documentUseCase.CreateNewVersion(c.Request.Context(), id, req.FilePath, req.FileSize, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, newVersion)
}

// GrantAccess grants access to a document
func (h *DocumentHandler) GrantAccess(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetUserID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.GrantAccess(c.Request.Context(), id, targetUserID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access granted successfully"})
}

// RevokeAccess revokes access to a document
func (h *DocumentHandler) RevokeAccess(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetUserID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.RevokeAccess(c.Request.Context(), id, targetUserID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access revoked successfully"})
}

// MakePublic makes a document public
func (h *DocumentHandler) MakePublic(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.MakePublic(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document is now public"})
}

// MakePrivate makes a document private
func (h *DocumentHandler) MakePrivate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.MakePrivate(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document is now private"})
}

// SetExpiration sets expiration date for a document
func (h *DocumentHandler) SetExpiration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req struct {
		ExpiresAt string `json:"expires_at" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use RFC3339 format"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.SetExpiration(c.Request.Context(), id, expiresAt, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expiration date set successfully"})
}

// GetDocumentsByEntity gets documents for an entity
func (h *DocumentHandler) GetDocumentsByEntity(c *gin.Context) {
	entityType := c.Param("entity_type")
	entityIDStr := c.Param("entity_id")

	entityID, err := primitive.ObjectIDFromHex(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	documents, err := h.documentUseCase.GetDocumentsByRelatedEntity(c.Request.Context(), entityType, entityID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents, "total": len(documents)})
}

// GetDocumentsByType gets documents by type
func (h *DocumentHandler) GetDocumentsByType(c *gin.Context) {
	docType := c.Param("type")
	userID := c.MustGet("user_id").(primitive.ObjectID)

	documents, err := h.documentUseCase.GetDocumentsByType(c.Request.Context(), entities.DocumentType(docType), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents, "total": len(documents)})
}

// GetDocumentsByCategory gets documents by category
func (h *DocumentHandler) GetDocumentsByCategory(c *gin.Context) {
	category := c.Param("category")
	userID := c.MustGet("user_id").(primitive.ObjectID)

	documents, err := h.documentUseCase.GetDocumentsByCategory(c.Request.Context(), entities.DocumentType(category), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents, "total": len(documents)})
}

// SearchDocuments searches documents
func (h *DocumentHandler) SearchDocuments(c *gin.Context) {
	filter := &repositories.DocumentFilter{}

	// Parse query parameters
	filter.Type = c.Query("type")
	filter.Search = c.Query("q")
	if filter.Search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}
	filter.IncludeArchived = c.Query("include_archived") == "true"

	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filter.Tags = tags
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "created_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	userID := c.MustGet("user_id").(primitive.ObjectID)

	documents, total, err := h.documentUseCase.SearchDocuments(c.Request.Context(), filter, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
		"total":     total,
		"limit":     filter.Limit,
		"offset":    filter.Offset,
	})
}

// ShareDocument shares a document
func (h *DocumentHandler) ShareDocument(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetUserID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.GrantAccess(c.Request.Context(), id, targetUserID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document shared successfully"})
}

// ArchiveDocument archives a document
func (h *DocumentHandler) ArchiveDocument(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.ArchiveDocument(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document archived successfully"})
}

// UnarchiveDocument unarchives a document
func (h *DocumentHandler) UnarchiveDocument(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.documentUseCase.UnarchiveDocument(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unarchived successfully"})
}
