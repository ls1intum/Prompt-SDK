package promptSDK

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ls1intum/prompt-sdk/promptTypes"
	log "github.com/sirupsen/logrus"
)

type Resolution struct {
	DtoName       string
	BaseURL       string
	EndpointPath  string
	CoursePhaseID uuid.UUID
}

type CoursePhaseParticipationsWithResolutions struct {
	Participations []promptTypes.CoursePhaseParticipationWithStudent `json:"participations"`
	Resolutions    []Resolution                                      `json:"resolutions"`
}

type PrevCoursePhaseData struct {
	PrevData    promptTypes.MetaData `json:"prevData"`
	Resolutions []Resolution         `json:"resolutions"`
}

// ────────────────────────────────────────────────────────────────────────────────
// URL helpers
// ────────────────────────────────────────────────────────────────────────────────

// transformBaseURL rewrites service URLs so that they point to the internal
// Docker network instead of the public host.  Extend the slice below to add
// more rules (pattern → replacement).
func transformBaseURL(baseURL string) string {
	rewriteRules := []struct {
		pattern     string
		replacement string
	}{
		// team‑allocation service
		{"/team-allocation/api", "http://server-team-allocation/team-allocation-api"},
		// intro‑course service
		{"/intro-course/api", "http://server-intro-course/intro-course/api"},
		{"/assessment/api", "http://server-assessment/assessment/api"},
	}

	for _, r := range rewriteRules {
		if strings.Contains(baseURL, r.pattern) {
			return r.replacement
		}
	}
	return baseURL // nothing matched – keep original
}

// buildURL constructs the request URL for a given resolution.
// extraPaths (such as a courseParticipationID) can be appended.
func buildURL(resolution Resolution, resolveLocally bool, extraPaths ...string) string {
	baseURL := resolution.BaseURL
	if resolveLocally {
		baseURL = transformBaseURL(baseURL)
	}

	url := fmt.Sprintf("%s/course_phase/%s/%s",
		baseURL,
		resolution.CoursePhaseID,
		getEndpointPath(resolution.EndpointPath),
	)

	if len(extraPaths) > 0 {
		url = fmt.Sprintf("%s/%s", url, strings.Join(extraPaths, "/"))
	}
	return url
}

// getEndpointPath trims leading and trailing slashes from the endpoint path.
func getEndpointPath(endpointPath string) string {
	return strings.Trim(endpointPath, "/")
}

// ────────────────────────────────────────────────────────────────────────────────
// Response helpers
// ────────────────────────────────────────────────────────────────────────────────

// parseAndValidate unmarshals the data into a map and ensures the expected key exists.
func parseAndValidate(data []byte, dtoName string) (interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	value, ok := result[dtoName]
	if !ok {
		log.Error("Failed to find expected key in response: ", dtoName)
		return nil, fmt.Errorf("failed to find expected key in response: %s", dtoName)
	}
	return value, nil
}

// ResolveParticipation resolves data for a single course participation.
func ResolveParticipation(authHeader string, resolution Resolution, courseParticipationID uuid.UUID, resolveLocally bool) (interface{}, error) {
	url := buildURL(resolution, resolveLocally, courseParticipationID.String())
	data, err := FetchJSON(url, authHeader)
	if err != nil {
		return nil, err
	}

	return parseAndValidate(data, resolution.DtoName)
}

// ResolveCoursePhaseData resolves data for a course phase.
func ResolveCoursePhaseData(authHeader string, resolution Resolution, resolveLocally bool) (interface{}, error) {
	url := buildURL(resolution, resolveLocally)
	data, err := FetchJSON(url, authHeader)
	if err != nil {
		return nil, err
	}

	return parseAndValidate(data, resolution.DtoName)
}

// ResolveAllParticipations resolves data for all participations and returns a map keyed by courseParticipationID.
func ResolveAllParticipations(authHeader string, resolution Resolution, resolveLocally bool) (map[uuid.UUID]interface{}, error) {
	url := buildURL(resolution, resolveLocally)
	data, err := FetchJSON(url, authHeader)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	formattedResults := make(map[uuid.UUID]interface{})
	for _, item := range results {
		idStr, ok := item["courseParticipationID"].(string)
		if !ok {
			log.Error("Failed to cast courseParticipationID to string")
			return nil, fmt.Errorf("failed to cast courseParticipationID to string")
		}
		participationID, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("Failed to parse courseParticipationID: ", err)
			return nil, fmt.Errorf("failed to parse courseParticipationID: %v", err)
		}
		formattedResults[participationID] = item[resolution.DtoName]
	}

	return formattedResults, nil
}

// FetchAndMergeParticipationsWithResolutions fetches participations and enriches each with resolved data.
func FetchAndMergeParticipationsWithResolutions(coreURL string, authHeader string, coursePhaseID uuid.UUID, resolveLocally bool) ([]promptTypes.CoursePhaseParticipationWithStudent, error) {

	url := fmt.Sprintf("%s/api/course_phases/%s/participations", coreURL, coursePhaseID)
	data, err := FetchJSON(url, authHeader)
	if err != nil {
		return nil, err
	}

	var payload CoursePhaseParticipationsWithResolutions
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	for _, res := range payload.Resolutions {
		resolved, err := ResolveAllParticipations(authHeader, res, resolveLocally)
		if err != nil {
			return nil, err
		}

		for i, p := range payload.Participations {
			if val, ok := resolved[p.CourseParticipationID]; ok {
				if p.PrevData == nil {
					p.PrevData = make(promptTypes.MetaData)
				}
				p.PrevData[res.DtoName] = val
				payload.Participations[i] = p
			}
		}
	}
	return payload.Participations, nil
}

func FetchAndMergeCoursePhaseWithResolution(coreURL string, authHeader string, coursePhaseID uuid.UUID, resolveLocally bool) (promptTypes.MetaData, error) {
	url := fmt.Sprintf("%s/api/course_phases/%s/course_phase_data", coreURL, coursePhaseID)
	data, err := FetchJSON(url, authHeader)
	if err != nil {
		return nil, err
	}

	var cpWithRes PrevCoursePhaseData
	if err := json.Unmarshal(data, &cpWithRes); err != nil {
		return nil, err
	}

	if cpWithRes.PrevData == nil {
		cpWithRes.PrevData = make(promptTypes.MetaData)
	}

	for _, res := range cpWithRes.Resolutions {
		resolvedData, err := ResolveCoursePhaseData(authHeader, res, resolveLocally)
		if err != nil {
			return nil, err
		}
		cpWithRes.PrevData[res.DtoName] = resolvedData
	}
	return cpWithRes.PrevData, nil
}
