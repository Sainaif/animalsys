package config

import "os"

type Config struct {
	MongoURI         string
	MongoDBName      string
	JWTSecret        string
	LDAPEnabled      bool
	LDAPServer       string
	LDAPBaseDN       string
	LDAPBindDN       string
	LDAPBindPassword string
	Port             string
}

func Load() Config {
	return Config{
		MongoURI:         getEnv("MONGO_URI", "mongodb://root:example@mongo:27017/animalsys?authSource=admin"),
		MongoDBName:      getEnv("MONGO_DB", "animalsys"),
		JWTSecret:        getEnv("JWT_SECRET", "supersecretkey"),
		LDAPEnabled:      getEnv("LDAP_ENABLED", "false") == "true",
		LDAPServer:       getEnv("LDAP_SERVER", ""),
		LDAPBaseDN:       getEnv("LDAP_BASE_DN", ""),
		LDAPBindDN:       getEnv("LDAP_BIND_DN", ""),
		LDAPBindPassword: getEnv("LDAP_BIND_PASSWORD", ""),
		Port:             getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
