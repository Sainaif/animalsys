package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRole represents the user's role in the system
type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleAdmin      UserRole = "admin"
	RoleEmployee   UserRole = "employee"
	RoleVolunteer  UserRole = "volunteer"
	RoleUser       UserRole = "user"
)

// UserStatus represents the user's account status
type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
)

// User represents a user in the system
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	FirstName    string             `bson:"first_name" json:"first_name" validate:"required"`
	LastName     string             `bson:"last_name" json:"last_name" validate:"required"`
	Role         UserRole           `bson:"role" json:"role" validate:"required"`
	Status       UserStatus         `bson:"status" json:"status" validate:"required"`
	Phone        string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Avatar       string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Language     string             `bson:"language" json:"language" validate:"required,oneof=en pl"`
	Theme        string             `bson:"theme" json:"theme" validate:"required,oneof=light dark"`
	RefreshToken string             `bson:"refresh_token,omitempty" json:"-"`
	LastLogin    *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// IsValidRole checks if the role is valid
func IsValidRole(role UserRole) bool {
	switch role {
	case RoleSuperAdmin, RoleAdmin, RoleEmployee, RoleVolunteer, RoleUser:
		return true
	}
	return false
}

// IsValidStatus checks if the status is valid
func IsValidStatus(status UserStatus) bool {
	switch status {
	case StatusActive, StatusInactive, StatusSuspended:
		return true
	}
	return false
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// IsSuperAdmin checks if the user is a super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// IsAdmin checks if the user is an admin or super admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// CanManageUsers checks if the user can manage other users
func (u *User) CanManageUsers() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAdmin
}

// CanEditSettings checks if the user can edit foundation settings
func (u *User) CanEditSettings() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAdmin
}
