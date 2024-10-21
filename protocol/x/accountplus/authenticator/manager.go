package authenticator

import (
	"sort"

	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

// AuthenticatorManager is a manager for all registered authenticators.
type AuthenticatorManager struct {
	registeredAuthenticators map[string]types.Authenticator
	orderedKeys              []string // slice to keep track of the keys in sorted order
}

// NewAuthenticatorManager creates a new AuthenticatorManager.
func NewAuthenticatorManager() *AuthenticatorManager {
	return &AuthenticatorManager{
		registeredAuthenticators: make(map[string]types.Authenticator),
		orderedKeys:              []string{},
	}
}

// ResetAuthenticators resets all registered authenticators.
func (am *AuthenticatorManager) ResetAuthenticators() {
	am.registeredAuthenticators = make(map[string]types.Authenticator)
	am.orderedKeys = []string{}
}

// InitializeAuthenticators initializes authenticators. If already initialized, it will not overwrite.
func (am *AuthenticatorManager) InitializeAuthenticators(initialAuthenticators []types.Authenticator) {
	if len(am.registeredAuthenticators) > 0 {
		return
	}
	for _, authenticator := range initialAuthenticators {
		am.registeredAuthenticators[authenticator.Type()] = authenticator
		am.orderedKeys = append(am.orderedKeys, authenticator.Type())
	}
	sort.Strings(am.orderedKeys) // Ensure keys are sorted
}

// RegisterAuthenticator adds a new authenticator to the map of registered authenticators.
func (am *AuthenticatorManager) RegisterAuthenticator(authenticator types.Authenticator) {
	if _, exists := am.registeredAuthenticators[authenticator.Type()]; !exists {
		am.orderedKeys = append(am.orderedKeys, authenticator.Type())
		sort.Strings(am.orderedKeys) // Re-sort keys after addition
	}
	am.registeredAuthenticators[authenticator.Type()] = authenticator
}

// UnregisterAuthenticator removes an authenticator from the map of registered authenticators.
func (am *AuthenticatorManager) UnregisterAuthenticator(authenticator types.Authenticator) {
	if _, exists := am.registeredAuthenticators[authenticator.Type()]; exists {
		delete(am.registeredAuthenticators, authenticator.Type())
		// Remove the key from orderedKeys
		for i, key := range am.orderedKeys {
			if key == authenticator.Type() {
				am.orderedKeys = append(am.orderedKeys[:i], am.orderedKeys[i+1:]...)
				break
			}
		}
	}
}

// GetRegisteredAuthenticators returns the list of registered authenticators in sorted order.
func (am *AuthenticatorManager) GetRegisteredAuthenticators() []types.Authenticator {
	var authenticators []types.Authenticator
	for _, key := range am.orderedKeys {
		authenticators = append(authenticators, am.registeredAuthenticators[key])
	}
	return authenticators
}

// IsAuthenticatorTypeRegistered checks if the authenticator type is registered.
func (am *AuthenticatorManager) IsAuthenticatorTypeRegistered(authenticatorType string) bool {
	_, exists := am.registeredAuthenticators[authenticatorType]
	return exists
}

// GetAuthenticatorByType returns the base implementation of the authenticator type.
func (am *AuthenticatorManager) GetAuthenticatorByType(authenticatorType string) types.Authenticator {
	if authenticator, exists := am.registeredAuthenticators[authenticatorType]; exists {
		return authenticator
	}
	return nil
}
