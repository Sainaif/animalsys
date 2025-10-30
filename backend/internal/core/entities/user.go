package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a system user
type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username          string             `bson:"username" json:"username"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bson:"password" json:"-"` // Never send password in JSON
	Role              string             `bson:"role" json:"role"`
	FirstName         string             `bson:"first_name" json:"first_name"`
	LastName          string             `bson:"last_name" json:"last_name"`
	Phone             string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Address           string             `bson:"address,omitempty" json:"address,omitempty"`
	ProfilePicture    string             `bson:"profile_picture,omitempty" json:"profile_picture,omitempty"`
	Active            bool               `bson:"active" json:"active"`
	EmailVerified     bool               `bson:"email_verified" json:"email_verified"`
	LastLogin         *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
	PasswordChangedAt *time.Time         `bson:"password_changed_at,omitempty" json:"password_changed_at,omitempty"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy         string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy         string             `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// UserCreateRequest represents user creation request
type UserCreateRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Role      string `json:"role" validate:"required,oneof=super_admin admin employee volunteer user"`
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

// UserUpdateRequest represents user update request
type UserUpdateRequest struct {
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
	Role      string `json:"role,omitempty" validate:"omitempty,oneof=super_admin admin employee volunteer user"`
	Active    *bool  `json:"active,omitempty"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UserFilter represents filters for querying users
type UserFilter struct {
	Role      string
	Active    *bool
	Search    string // Search in username, email, first name, last name
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string // asc or desc
}

// NewUser creates a new user
func NewUser(username, email, hashedPassword, role, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		ID:            primitive.NewObjectID(),
		Username:      username,
		Email:         email,
		Password:      hashedPassword,
		Role:          role,
		FirstName:     firstName,
		LastName:      lastName,
		Active:        true,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// FullName returns user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

// UpdatePassword updates the password and password change timestamp
func (u *User) UpdatePassword(hashedPassword string) {
	now := time.Now()
	u.Password = hashedPassword
	u.PasswordChangedAt = &now
	u.UpdatedAt = now
}

// Deactivate deactivates the user
func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user
func (u *User) Activate() {
	u.Active = true
	u.UpdatedAt = time.Now()
}

// VerifyEmail marks email as verified
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now()
}
