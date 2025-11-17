package main

import (
	"fmt"
	"os"
	"strings"
)

type organizationConfig struct {
	Name                string
	ShortName           string
	LegalName           string
	Domain              string
	Description         string
	Email               string
	Phone               string
	Website             string
	Address             string
	CenterName          string
	DefaultPassword     string
	SuperAdminFirstName string
	SuperAdminLastName  string
	SuperAdminEmail     string
}

func loadOrganizationConfig() organizationConfig {
	name := getEnv("ORG_NAME", "Happy Paws Animal Foundation")
	shortName := getEnv("ORG_SHORT_NAME", "Animal Foundation")
	legalName := getEnv("ORG_LEGAL_NAME", name)
	domain := strings.ToLower(getEnv("ORG_DOMAIN", "happypaws.org"))
	center := getEnv("ORG_CENTER_NAME", fmt.Sprintf("%s Center", shortName))

	cfg := organizationConfig{
		Name:                name,
		ShortName:           shortName,
		LegalName:           legalName,
		Domain:              domain,
		Description:         getEnv("ORG_DESCRIPTION", "Dedicated to rescuing and rehoming animals in need"),
		Email:               getEnv("ORG_EMAIL", fmt.Sprintf("info@%s", domain)),
		Phone:               getEnv("ORG_PHONE", "+1-555-ANIMALS"),
		Website:             getEnv("ORG_WEBSITE", fmt.Sprintf("https://%s", strings.TrimPrefix(domain, "www."))),
		Address:             getEnv("ORG_ADDRESS", "123 Rescue Lane, Animal City, AC 12345"),
		CenterName:          center,
		DefaultPassword:     getEnv("ORG_DEFAULT_PASSWORD", "password123"),
		SuperAdminFirstName: getEnv("ORG_SUPER_ADMIN_FIRST_NAME", "Sarah"),
		SuperAdminLastName:  getEnv("ORG_SUPER_ADMIN_LAST_NAME", "Johnson"),
	}

	defaultSuperAdminEmail := fmt.Sprintf("%s.%s@%s",
		strings.ToLower(cfg.SuperAdminFirstName),
		strings.ToLower(cfg.SuperAdminLastName),
		cfg.Domain,
	)
	cfg.SuperAdminEmail = getEnv("ORG_SUPER_ADMIN_EMAIL", defaultSuperAdminEmail)

	return cfg
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
