package keycloakCoreRequests

import (
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func sendRequest(coreURL, method, subURL, authHeader string, body io.Reader) (*http.Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	requestURL := coreURL + subURL
	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		log.Error("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error sending request:", err)
		return nil, err
	}

	return resp, nil
}
