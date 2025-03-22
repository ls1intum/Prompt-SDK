package promptSDK

import (
	"github.com/gin-gonic/gin"
	"github.com/ls1intum/prompt-sdk/keycloakTokenVerifier"
)

// exposing the roles
const PromptAdmin = keycloakTokenVerifier.PromptAdmin
const PromptLecturer = keycloakTokenVerifier.PromptLecturer
const CourseLecturer = keycloakTokenVerifier.CourseLecturer
const CourseEditor = keycloakTokenVerifier.CourseEditor
const CourseStudent = keycloakTokenVerifier.CourseStudent

func InitAuthenticationMiddleware(KeycloakURL, Realm, CoreURL string) error {
	return keycloakTokenVerifier.InitKeycloakTokenVerifier(KeycloakURL, Realm, CoreURL)
}

func AuthenticationMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return keycloakTokenVerifier.AuthenticationMiddleware(allowedRoles...)
}
