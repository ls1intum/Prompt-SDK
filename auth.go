package promptSDK

import (
	"github.com/gin-gonic/gin"
	"github.com/ls1intum/prompt-sdk/keycloakTokenVerifier"
)

func InitAuthenticationMiddleware(KeycloakURL, Realm, CoreURL string) error {
	return keycloakTokenVerifier.InitKeycloakTokenVerifier(KeycloakURL, Realm, CoreURL)
}

func AuthenticationMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return keycloakTokenVerifier.AuthenticationMiddleware(allowedRoles...)
}
