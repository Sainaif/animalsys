package utils

import (
	"animalsys/config"
	"errors"

	ldap "github.com/go-ldap/ldap/v3"
)

func LDAPAuthenticate(username, password string, cfg config.Config) (string, error) {
	if !cfg.LDAPEnabled {
		return "", errors.New("LDAP not enabled")
	}
	conn, err := ldap.DialURL(cfg.LDAPServer)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	bindDN := cfg.LDAPBindDN
	bindPass := cfg.LDAPBindPassword
	if err := conn.Bind(bindDN, bindPass); err != nil {
		return "", err
	}
	searchRequest := ldap.NewSearchRequest(
		cfg.LDAPBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=person)(uid="+ldap.EscapeFilter(username)+"))",
		[]string{"dn"},
		nil,
	)
	result, err := conn.Search(searchRequest)
	if err != nil || len(result.Entries) == 0 {
		return "", errors.New("LDAP user not found")
	}
	userDN := result.Entries[0].DN
	if err := conn.Bind(userDN, password); err != nil {
		return "", errors.New("LDAP auth failed")
	}
	return userDN, nil
}
