package keycloakTokenVerifier

import "github.com/sirupsen/logrus"

type KeycloakTokenVerifier struct {
	KeycloakURL             string
	Realm                   string
	ClientID                string
	expectedAuthorizedParty string
	CoreURL                 string
}

var KeycloakTokenVerifierSingleton *KeycloakTokenVerifier

func InitKeycloakTokenVerifier(KeycloakURL, Realm, CoreURL string) error {
	KeycloakTokenVerifierSingleton = &KeycloakTokenVerifier{
		KeycloakURL:             KeycloakURL,
		Realm:                   Realm,
		ClientID:                "prompt-server",
		expectedAuthorizedParty: "prompt-client",
		CoreURL:                 CoreURL,
	}

	// init the middleware
	err := InitKeycloakVerifier()
	if err != nil {
		logrus.Error("Failed to initialize keycloak verifier: ", err)
		return err
	}
	return nil
}
