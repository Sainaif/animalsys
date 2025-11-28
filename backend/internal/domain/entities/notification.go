package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeSuccess NotificationType = "success"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeError   NotificationType = "error"
)

// NotificationPriority represents the priority of a notification
type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityNormal NotificationPriority = "normal"
	NotificationPriorityHigh   NotificationPriority = "high"
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

// Notification represents an in-app notification
type Notification struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Recipient
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`

	// Notification Details
	Type     NotificationType     `json:"type" bson:"type"`
	Priority NotificationPriority `json:"priority" bson:"priority"`
	Title    string               `json:"title" bson:"title"`
	Message  string               `json:"message" bson:"message"`
	Icon     string               `json:"icon,omitempty" bson:"icon,omitempty"`

	// Status
	Read      bool       `json:"read" bson:"read"`
	ReadAt    *time.Time `json:"read_at,omitempty" bson:"read_at,omitempty"`
	Dismissed bool       `json:"dismissed" bson:"dismissed"`

	// Action/Link
	ActionURL    string `json:"action_url,omitempty" bson:"action_url,omitempty"`
	ActionText   string `json:"action_text,omitempty" bson:"action_text,omitempty"`

	// Related Resource
	RelatedType string              `json:"related_type,omitempty" bson:"related_type,omitempty"` // animal, adoption, event, etc.
	RelatedID   *primitive.ObjectID `json:"related_id,omitempty" bson:"related_id,omitempty"`

	// Grouping
	Category string `json:"category,omitempty" bson:"category,omitempty"` // For grouping similar notifications
	GroupKey string `json:"group_key,omitempty" bson:"group_key,omitempty"` // For replacing/updating similar notifications

	// Expiration
	ExpiresAt *time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`

	// Additional Data
	Metadata map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	if !n.Read {
		now := time.Now()
		n.Read = true
		n.ReadAt = &now
		n.UpdatedAt = now
	}
}

// MarkAsUnread marks the notification as unread
func (n *Notification) MarkAsUnread() {
	if n.Read {
		n.Read = false
		n.ReadAt = nil
		n.UpdatedAt = time.Now()
	}
}

// Dismiss dismisses the notification
func (n *Notification) Dismiss() {
	n.Dismissed = true
	n.UpdatedAt = time.Now()
}

// IsExpired checks if the notification has expired
func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return n.ExpiresAt.Before(time.Now())
}

// NewNotification creates a new notification
func NewNotification(
	userID primitive.ObjectID,
	notifType NotificationType,
	title string,
	message string,
) *Notification {
	now := time.Now()
	return &Notification{
		UserID:    userID,
		Type:      notifType,
		Priority:  NotificationPriorityNormal,
		Title:     title,
		Message:   message,
		Read:      false,
		Dismissed: false,
		Metadata:  make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewInfoNotification creates a new info notification
func NewInfoNotification(userID primitive.ObjectID, title, message string) *Notification {
	return NewNotification(userID, NotificationTypeInfo, title, message)
}

// NewSuccessNotification creates a new success notification
func NewSuccessNotification(userID primitive.ObjectID, title, message string) *Notification {
	notif := NewNotification(userID, NotificationTypeSuccess, title, message)
	notif.Icon = "check-circle"
	return notif
}

// NewWarningNotification creates a new warning notification
func NewWarningNotification(userID primitive.ObjectID, title, message string) *Notification {
	notif := NewNotification(userID, NotificationTypeWarning, title, message)
	notif.Priority = NotificationPriorityHigh
	notif.Icon = "exclamation-triangle"
	return notif
}

// NewErrorNotification creates a new error notification
func NewErrorNotification(userID primitive.ObjectID, title, message string) *Notification {
	notif := NewNotification(userID, NotificationTypeError, title, message)
	notif.Priority = NotificationPriorityUrgent
	notif.Icon = "exclamation-circle"
	return notif
}

// NotificationPreferences represents user-specific notification settings
type NotificationPreferences struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	EmailEnabled bool               `json:"email_enabled" bson:"email_enabled"`
	PushEnabled  bool               `json:"push_enabled" bson:"push_enabled"`
	Categories   map[string]bool    `json:"categories" bson:"categories"` // e.g., "adoption", "events"
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewNotificationPreferences creates default preferences for a user
func NewNotificationPreferences(userID primitive.ObjectID) *NotificationPreferences {
	return &NotificationPreferences{
		UserID:       userID,
		EmailEnabled: true,
		PushEnabled:  false,
		Categories:   make(map[string]bool), // Default to all categories enabled (or handle logic in use case)
		UpdatedAt:    time.Now(),
	}
}
