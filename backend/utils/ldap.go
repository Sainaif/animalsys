package utils

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type LDAPConfig struct {
	Server       string
	BaseDN       string
	BindDN       string
	BindPassword string
}

func AuthenticateLDAP(username, password string, config LDAPConfig) (string, error) {
	l, err := ldap.DialURL(config.Server)
	if err != nil {
		return "", fmt.Errorf("failed to connect to LDAP server: %w", err)
	}
	defer l.Close()

	err = l.Bind(config.BindDN, config.BindPassword)
	if err != nil {
		return "", fmt.Errorf("failed to bind to LDAP: %w", err)
	}

	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(uid=%s)", username),
		[]string{"dn", "cn", "mail"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return "", fmt.Errorf("failed to search LDAP: %w", err)
	}

	if len(sr.Entries) != 1 {
		return "", fmt.Errorf("user not found or too many entries")
	}

	userDN := sr.Entries[0].DN

	err = l.Bind(userDN, password)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	return userDN, nil
}
