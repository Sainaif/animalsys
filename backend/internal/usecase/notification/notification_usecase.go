package notification

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationUseCase struct {
	notificationRepo repositories.NotificationRepository
	auditLogRepo     repositories.AuditLogRepository
}

func NewNotificationUseCase(
	notificationRepo repositories.NotificationRepository,
	auditLogRepo repositories.AuditLogRepository,
) *NotificationUseCase {
	return &NotificationUseCase{
		notificationRepo: notificationRepo,
		auditLogRepo:     auditLogRepo,
	}
}

// CreateNotification creates a new notification
func (uc *NotificationUseCase) CreateNotification(ctx context.Context, notification *entities.Notification) error {
	// Validate required fields
	if notification.Title == "" {
		return errors.NewBadRequest("Notification title is required")
	}

	if notification.Message == "" {
		return errors.NewBadRequest("Notification message is required")
	}

	// Check for existing notification with same group key
	if notification.GroupKey != "" {
		existing, err := uc.notificationRepo.FindByGroupKey(ctx, notification.UserID, notification.GroupKey)
		if err == nil && existing != nil {
			// Update existing notification instead of creating new one
			existing.Title = notification.Title
			existing.Message = notification.Message
			existing.Type = notification.Type
			existing.Priority = notification.Priority
			existing.ActionURL = notification.ActionURL
			existing.ActionText = notification.ActionText
			existing.Icon = notification.Icon
			existing.Metadata = notification.Metadata
			existing.Read = false // Mark as unread
			existing.ReadAt = nil
			return uc.notificationRepo.Update(ctx, existing)
		}
	}

	return uc.notificationRepo.Create(ctx, notification)
}

// GetNotificationByID retrieves a notification by ID
func (uc *NotificationUseCase) GetNotificationByID(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*entities.Notification, error) {
	notification, err := uc.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify the notification belongs to the user
	if notification.UserID != userID {
		return nil, errors.NewForbidden("You don't have permission to access this notification")
	}

	return notification, nil
}

// ListNotifications lists notifications with filtering
func (uc *NotificationUseCase) ListNotifications(ctx context.Context, filter *repositories.NotificationFilter) ([]*entities.Notification, int64, error) {
	return uc.notificationRepo.List(ctx, filter)
}

// GetUserNotifications retrieves all non-dismissed notifications for a user
func (uc *NotificationUseCase) GetUserNotifications(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error) {
	return uc.notificationRepo.GetByUser(ctx, userID)
}

// GetUnreadNotifications retrieves unread notifications for a user
func (uc *NotificationUseCase) GetUnreadNotifications(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error) {
	return uc.notificationRepo.GetUnreadByUser(ctx, userID)
}

// GetUnreadCount retrieves the count of unread notifications for a user
func (uc *NotificationUseCase) GetUnreadCount(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return uc.notificationRepo.CountUnreadByUser(ctx, userID)
}

// MarkAsRead marks a notification as read
func (uc *NotificationUseCase) MarkAsRead(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	notification, err := uc.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify the notification belongs to the user
	if notification.UserID != userID {
		return errors.NewForbidden("You don't have permission to modify this notification")
	}

	if notification.Read {
		return nil // Already read
	}

	if err := uc.notificationRepo.MarkAsRead(ctx, id); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "notification", "", "").
			WithEntityID(notification.ID).
			WithChanges(map[string]interface{}{
				"read": true,
			}))

	return nil
}

// MarkAsUnread marks a notification as unread
func (uc *NotificationUseCase) MarkAsUnread(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	notification, err := uc.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify the notification belongs to the user
	if notification.UserID != userID {
		return errors.NewForbidden("You don't have permission to modify this notification")
	}

	if !notification.Read {
		return nil // Already unread
	}

	if err := uc.notificationRepo.MarkAsUnread(ctx, id); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "notification", "", "").
			WithEntityID(notification.ID).
			WithChanges(map[string]interface{}{
				"read": false,
			}))

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (uc *NotificationUseCase) MarkAllAsRead(ctx context.Context, userID primitive.ObjectID) error {
	if err := uc.notificationRepo.MarkAllAsRead(ctx, userID); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "notification", "", "").
			WithChanges(map[string]interface{}{
				"action": "mark_all_as_read",
			}))

	return nil
}

// DismissNotification dismisses a notification
func (uc *NotificationUseCase) DismissNotification(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	notification, err := uc.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify the notification belongs to the user
	if notification.UserID != userID {
		return errors.NewForbidden("You don't have permission to modify this notification")
	}

	if notification.Dismissed {
		return nil // Already dismissed
	}

	if err := uc.notificationRepo.Dismiss(ctx, id); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "notification", "", "").
			WithEntityID(notification.ID).
			WithChanges(map[string]interface{}{
				"dismissed": true,
			}))

	return nil
}

// DeleteNotification deletes a notification
func (uc *NotificationUseCase) DeleteNotification(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	notification, err := uc.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify the notification belongs to the user
	if notification.UserID != userID {
		return errors.NewForbidden("You don't have permission to delete this notification")
	}

	if err := uc.notificationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "notification", "", "").
			WithEntityID(notification.ID))

	return nil
}

// DeleteExpiredNotifications deletes all expired notifications
func (uc *NotificationUseCase) DeleteExpiredNotifications(ctx context.Context) (int64, error) {
	return uc.notificationRepo.DeleteExpired(ctx)
}

// NotifyUser creates a notification for a user
func (uc *NotificationUseCase) NotifyUser(
	ctx context.Context,
	userID primitive.ObjectID,
	notifType entities.NotificationType,
	title string,
	message string,
) error {
	notification := entities.NewNotification(userID, notifType, title, message)
	return uc.CreateNotification(ctx, notification)
}

// NotifyUserWithAction creates a notification with an action button
func (uc *NotificationUseCase) NotifyUserWithAction(
	ctx context.Context,
	userID primitive.ObjectID,
	notifType entities.NotificationType,
	title string,
	message string,
	actionURL string,
	actionText string,
) error {
	notification := entities.NewNotification(userID, notifType, title, message)
	notification.ActionURL = actionURL
	notification.ActionText = actionText
	return uc.CreateNotification(ctx, notification)
}

// NotifyUserWithRelated creates a notification related to another resource
func (uc *NotificationUseCase) NotifyUserWithRelated(
	ctx context.Context,
	userID primitive.ObjectID,
	notifType entities.NotificationType,
	title string,
	message string,
	relatedType string,
	relatedID primitive.ObjectID,
) error {
	notification := entities.NewNotification(userID, notifType, title, message)
	notification.RelatedType = relatedType
	notification.RelatedID = &relatedID
	return uc.CreateNotification(ctx, notification)
}

// NotifySuccess creates a success notification
func (uc *NotificationUseCase) NotifySuccess(ctx context.Context, userID primitive.ObjectID, title, message string) error {
	notification := entities.NewSuccessNotification(userID, title, message)
	return uc.CreateNotification(ctx, notification)
}

// NotifyWarning creates a warning notification
func (uc *NotificationUseCase) NotifyWarning(ctx context.Context, userID primitive.ObjectID, title, message string) error {
	notification := entities.NewWarningNotification(userID, title, message)
	return uc.CreateNotification(ctx, notification)
}

// NotifyError creates an error notification
func (uc *NotificationUseCase) NotifyError(ctx context.Context, userID primitive.ObjectID, title, message string) error {
	notification := entities.NewErrorNotification(userID, title, message)
	return uc.CreateNotification(ctx, notification)
}
