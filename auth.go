package promptSDK

import (
	"github.com/ls1intum/prompt-sdk/keycloakTokenVerifier"
)

func InitAuthenticationMiddleware(KeycloakURL, Realm, CoreURL string) error {
	return keycloakTokenVerifier.InitKeycloakTokenVerifier(KeycloakURL, Realm, CoreURL)
}
