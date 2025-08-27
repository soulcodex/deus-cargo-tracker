package postgres

import (
	"fmt"
)

type SSLMode string

func (s SSLMode) String() string {
	return string(s)
}

func (s SSLMode) validate() error {
	if _, ok := s.allowed()[s.String()]; !ok {
		return fmt.Errorf("invalid SSL mode: %s", s.String())
	}

	return nil
}

func (s SSLMode) allowed() map[string]SSLMode {
	return map[string]SSLMode{
		"disable":     DisableMode,
		"allow":       AllowMode,
		"prefer":      PreferMode,
		"require":     RequireMode,
		"verify-ca":   VerifyCAMode,
		"verify-full": VerifyFullMode,
	}
}

func SSLModeFromString(mode string) (SSLMode, error) {
	sslMode := SSLMode(mode)

	if err := sslMode.validate(); err != nil {
		return "", err
	}

	return sslMode, nil
}

const (
	DisableMode    SSLMode = "disable"
	AllowMode      SSLMode = "allow"
	PreferMode     SSLMode = "prefer"
	RequireMode    SSLMode = "require"
	VerifyCAMode   SSLMode = "verify-ca"
	VerifyFullMode SSLMode = "verify-full"
)
