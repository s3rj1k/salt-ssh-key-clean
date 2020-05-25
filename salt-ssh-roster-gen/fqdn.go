package main

import (
	"strings"

	"golang.org/x/net/publicsuffix"
)

// GetShortFQDN returns effective TLD plus one and handles domains without subdomains (plus one).
func GetShortFQDN(fqdn string) string {
	var (
		suff  string
		icann bool
		err   error
	)

	// minimal sanitization
	fqdn = strings.ToLower(strings.TrimSpace(fqdn))

	suff, err = publicsuffix.EffectiveTLDPlusOne(fqdn)
	if err != nil {
		return fqdn
	}

	if fqdn == suff { // no GTLD+1
		suff, icann = publicsuffix.PublicSuffix(fqdn)
		if !icann { // not registered in ICANN
			return fqdn
		}
	}

	return strings.TrimSuffix(
		strings.TrimSuffix(fqdn, suff),
		".",
	)
}

// GetFQDNWithOutPublicSuffix returns domain stripped of public suffix (ICANN registered).
func GetFQDNWithOutPublicSuffix(fqdn string) string {
	// minimal sanitization
	fqdn = strings.ToLower(strings.TrimSpace(fqdn))

	suff, icann := publicsuffix.PublicSuffix(fqdn)
	if !icann { // not registered in ICANN
		return fqdn
	}

	return strings.TrimSuffix(
		strings.TrimSuffix(fqdn, suff),
		".",
	)
}
