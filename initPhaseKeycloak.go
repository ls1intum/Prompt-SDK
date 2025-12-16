package promptSDK

import (
	"strings"

	"github.com/ls1intum/prompt-sdk/utils"
	log "github.com/sirupsen/logrus"
)

//nolint:unused // Public SDK function for external use
func initPhaseKeycloak() {
	baseURL := GetEnv("KEYCLOAK_HOST", "http://localhost:8081")
	if !strings.HasPrefix(baseURL, "http") {
		log.Warn("Keycloak host does not start with http(s). Adding https:// as prefix.")
		baseURL = "https://" + baseURL
	}

	realm := GetEnv("KEYCLOAK_REALM_NAME", "prompt")

	coreURL := utils.GetCoreUrl()
	err := InitAuthenticationMiddleware(baseURL, realm, coreURL)
	if err != nil {
		log.Fatalf("Failed to initialize keycloak: %v", err)
	}
}
