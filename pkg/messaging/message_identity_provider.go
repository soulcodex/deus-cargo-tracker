package messaging

import (
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

var (
	defaultIdentityProvider = UUIDIdentityProvider()
)

// SetIdentityProvider sets the identity provider for the messaging package.
func SetIdentityProvider(provider IdentityProvider) {
	defaultIdentityProvider = provider
}

// IdentityProvider is an interface for providing identity information.
type IdentityProvider interface {
	// Provide returns a new unique resource identifier.
	Provide() string
}

type identityProviderWrapper struct {
	provider func() string
}

func (w identityProviderWrapper) Provide() string {
	return w.provider()
}

// UUIDIdentityProvider returns an identity provider that generates UUIDs.
func UUIDIdentityProvider() IdentityProvider {
	return wrapUUIDIdentityProvider(utils.NewRandomUUIDProvider())
}

func wrapUUIDIdentityProvider(provider utils.UUIDProvider) IdentityProvider {
	return identityProviderWrapper{
		provider: func() string {
			return provider.New().String()
		},
	}
}
