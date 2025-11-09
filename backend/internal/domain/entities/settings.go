package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FoundationSettings represents the global foundation settings
type FoundationSettings struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Foundation Information
	Name             string            `json:"name" bson:"name"`
	LegalName        string            `json:"legal_name,omitempty" bson:"legal_name,omitempty"`
	Description      string            `json:"description,omitempty" bson:"description,omitempty"`
	Mission          string            `json:"mission,omitempty" bson:"mission,omitempty"`
	Vision           string            `json:"vision,omitempty" bson:"vision,omitempty"`
	FoundedDate      *time.Time        `json:"founded_date,omitempty" bson:"founded_date,omitempty"`
	TaxID            string            `json:"tax_id,omitempty" bson:"tax_id,omitempty"`
	RegistrationNumber string          `json:"registration_number,omitempty" bson:"registration_number,omitempty"`

	// Contact Information
	ContactInfo ContactDetails `json:"contact_info" bson:"contact_info"`

	// Address
	Address AddressInfo `json:"address,omitempty" bson:"address,omitempty"`

	// Social Media
	SocialMedia map[string]string `json:"social_media,omitempty" bson:"social_media,omitempty"` // platform -> URL

	// Operating Hours
	OperatingHours map[string]OperatingHour `json:"operating_hours,omitempty" bson:"operating_hours,omitempty"` // day -> hours

	// Policies
	AdoptionPolicy     string `json:"adoption_policy,omitempty" bson:"adoption_policy,omitempty"`
	PrivacyPolicy      string `json:"privacy_policy,omitempty" bson:"privacy_policy,omitempty"`
	TermsOfService     string `json:"terms_of_service,omitempty" bson:"terms_of_service,omitempty"`
	VolunteerPolicy    string `json:"volunteer_policy,omitempty" bson:"volunteer_policy,omitempty"`
	DonationPolicy     string `json:"donation_policy,omitempty" bson:"donation_policy,omitempty"`

	// Fees & Pricing
	DefaultAdoptionFees map[string]float64 `json:"default_adoption_fees,omitempty" bson:"default_adoption_fees,omitempty"` // species -> fee

	// Email Settings
	EmailSettings EmailSettings `json:"email_settings" bson:"email_settings"`

	// Notification Settings
	NotificationSettings NotificationSettings `json:"notification_settings" bson:"notification_settings"`

	// Features
	Features FeatureFlags `json:"features" bson:"features"`

	// Limits
	Limits SystemLimits `json:"limits" bson:"limits"`

	// Customization
	Branding Branding `json:"branding,omitempty" bson:"branding,omitempty"`

	// Metadata
	Version   int64      `json:"version" bson:"version"` // For optimistic locking
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
}

// ContactDetails represents contact information
type ContactDetails struct {
	Email         string   `json:"email" bson:"email"`
	Phone         string   `json:"phone,omitempty" bson:"phone,omitempty"`
	Fax           string   `json:"fax,omitempty" bson:"fax,omitempty"`
	Mobile        string   `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Website       string   `json:"website,omitempty" bson:"website,omitempty"`
	EmergencyPhone string  `json:"emergency_phone,omitempty" bson:"emergency_phone,omitempty"`
}

// OperatingHour represents operating hours for a day
type OperatingHour struct {
	Open     string `json:"open" bson:"open"`           // e.g., "09:00"
	Close    string `json:"close" bson:"close"`         // e.g., "17:00"
	IsClosed bool   `json:"is_closed" bson:"is_closed"` // true if closed that day
}

// EmailSettings represents email configuration
type EmailSettings struct {
	SMTPHost       string `json:"smtp_host,omitempty" bson:"smtp_host,omitempty"`
	SMTPPort       int    `json:"smtp_port,omitempty" bson:"smtp_port,omitempty"`
	SMTPUsername   string `json:"smtp_username,omitempty" bson:"smtp_username,omitempty"`
	SMTPPassword   string `json:"-" bson:"smtp_password,omitempty"` // Hidden from JSON
	FromEmail      string `json:"from_email" bson:"from_email"`
	FromName       string `json:"from_name" bson:"from_name"`
	ReplyToEmail   string `json:"reply_to_email,omitempty" bson:"reply_to_email,omitempty"`
	EnableTLS      bool   `json:"enable_tls" bson:"enable_tls"`
	EmailSignature string `json:"email_signature,omitempty" bson:"email_signature,omitempty"`
}

// NotificationSettings represents notification preferences
type NotificationSettings struct {
	EnableEmailNotifications   bool `json:"enable_email_notifications" bson:"enable_email_notifications"`
	EnableSMSNotifications     bool `json:"enable_sms_notifications" bson:"enable_sms_notifications"`
	EnablePushNotifications    bool `json:"enable_push_notifications" bson:"enable_push_notifications"`
	NotifyOnNewAdoption        bool `json:"notify_on_new_adoption" bson:"notify_on_new_adoption"`
	NotifyOnNewDonation        bool `json:"notify_on_new_donation" bson:"notify_on_new_donation"`
	NotifyOnNewVolunteer       bool `json:"notify_on_new_volunteer" bson:"notify_on_new_volunteer"`
	NotifyOnAnimalIntake       bool `json:"notify_on_animal_intake" bson:"notify_on_animal_intake"`
	NotifyOnVeterinaryVisit    bool `json:"notify_on_veterinary_visit" bson:"notify_on_veterinary_visit"`
	NotifyOnLowInventory       bool `json:"notify_on_low_inventory" bson:"notify_on_low_inventory"`
	DigestFrequency            string `json:"digest_frequency,omitempty" bson:"digest_frequency,omitempty"` // daily, weekly, never
}

// FeatureFlags controls which features are enabled
type FeatureFlags struct {
	EnableAdoptions       bool `json:"enable_adoptions" bson:"enable_adoptions"`
	EnableDonations       bool `json:"enable_donations" bson:"enable_donations"`
	EnableVolunteers      bool `json:"enable_volunteers" bson:"enable_volunteers"`
	EnableEvents          bool `json:"enable_events" bson:"enable_events"`
	EnableCampaigns       bool `json:"enable_campaigns" bson:"enable_campaigns"`
	EnableReports         bool `json:"enable_reports" bson:"enable_reports"`
	EnablePublicAPI       bool `json:"enable_public_api" bson:"enable_public_api"`
	EnableOnlineAdoption  bool `json:"enable_online_adoption" bson:"enable_online_adoption"`
	EnableOnlineDonation  bool `json:"enable_online_donation" bson:"enable_online_donation"`
	EnableNewsletters     bool `json:"enable_newsletters" bson:"enable_newsletters"`
	EnableSocialSharing   bool `json:"enable_social_sharing" bson:"enable_social_sharing"`
	MaintenanceMode       bool `json:"maintenance_mode" bson:"maintenance_mode"`
}

// SystemLimits defines various system limits
type SystemLimits struct {
	MaxAnimalsPerPage       int64 `json:"max_animals_per_page" bson:"max_animals_per_page"`
	MaxDonationsPerPage     int64 `json:"max_donations_per_page" bson:"max_donations_per_page"`
	MaxVolunteersPerPage    int64 `json:"max_volunteers_per_page" bson:"max_volunteers_per_page"`
	MaxFileUploadSizeMB     int64 `json:"max_file_upload_size_mb" bson:"max_file_upload_size_mb"`
	MaxImageUploadSizeMB    int64 `json:"max_image_upload_size_mb" bson:"max_image_upload_size_mb"`
	MaxImagesPerAnimal      int64 `json:"max_images_per_animal" bson:"max_images_per_animal"`
	SessionTimeoutMinutes   int64 `json:"session_timeout_minutes" bson:"session_timeout_minutes"`
	PasswordExpiryDays      int64 `json:"password_expiry_days" bson:"password_expiry_days"`
	MaxLoginAttempts        int64 `json:"max_login_attempts" bson:"max_login_attempts"`
	DataRetentionDays       int64 `json:"data_retention_days" bson:"data_retention_days"`
}

// Branding represents visual customization
type Branding struct {
	LogoURL          string            `json:"logo_url,omitempty" bson:"logo_url,omitempty"`
	FaviconURL       string            `json:"favicon_url,omitempty" bson:"favicon_url,omitempty"`
	PrimaryColor     string            `json:"primary_color,omitempty" bson:"primary_color,omitempty"`
	SecondaryColor   string            `json:"secondary_color,omitempty" bson:"secondary_color,omitempty"`
	AccentColor      string            `json:"accent_color,omitempty" bson:"accent_color,omitempty"`
	CustomCSS        string            `json:"custom_css,omitempty" bson:"custom_css,omitempty"`
	CustomFooter     string            `json:"custom_footer,omitempty" bson:"custom_footer,omitempty"`
	CustomHeader     string            `json:"custom_header,omitempty" bson:"custom_header,omitempty"`
}

// NewFoundationSettings creates a new foundation settings with defaults
func NewFoundationSettings(name string, updatedBy primitive.ObjectID) *FoundationSettings {
	now := time.Now()
	return &FoundationSettings{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Version:   1,
		UpdatedBy: updatedBy,
		UpdatedAt: now,
		CreatedAt: now,
		ContactInfo: ContactDetails{},
		EmailSettings: EmailSettings{
			EnableTLS: true,
			SMTPPort:  587,
		},
		NotificationSettings: NotificationSettings{
			EnableEmailNotifications: true,
			NotifyOnNewAdoption:     true,
			NotifyOnNewDonation:     true,
			DigestFrequency:         "daily",
		},
		Features: FeatureFlags{
			EnableAdoptions:      true,
			EnableDonations:      true,
			EnableVolunteers:     true,
			EnableEvents:         true,
			EnableCampaigns:      true,
			EnableReports:        true,
			EnableOnlineAdoption: true,
			EnableOnlineDonation: true,
			MaintenanceMode:      false,
		},
		Limits: SystemLimits{
			MaxAnimalsPerPage:     50,
			MaxDonationsPerPage:   50,
			MaxVolunteersPerPage:  50,
			MaxFileUploadSizeMB:   10,
			MaxImageUploadSizeMB:  5,
			MaxImagesPerAnimal:    10,
			SessionTimeoutMinutes: 30,
			PasswordExpiryDays:    90,
			MaxLoginAttempts:      5,
			DataRetentionDays:     365,
		},
		DefaultAdoptionFees: make(map[string]float64),
		SocialMedia:        make(map[string]string),
		OperatingHours:     make(map[string]OperatingHour),
	}
}

// UpdateVersion increments the version for optimistic locking
func (s *FoundationSettings) UpdateVersion(updatedBy primitive.ObjectID) {
	s.Version++
	s.UpdatedBy = updatedBy
	s.UpdatedAt = time.Now()
}
