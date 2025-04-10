package promptSDK

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

type Resolution struct {
	DtoName       string
	BaseURL       string
	EndpointPath  string
	CoursePhaseID uuid.UUID
}

type MetaData map[string]interface{}

type Student struct {
	ID                   uuid.UUID   `json:"id"`
	FirstName            string      `json:"firstName"`
	LastName             string      `json:"lastName"`
	Email                string      `json:"email"`
	MatriculationNumber  string      `json:"matriculationNumber"`
	UniversityLogin      string      `json:"universityLogin"`
	HasUniversityAccount bool        `json:"hasUniversityAccount"`
	Gender               string      `json:"gender"` // for simplicity we map the enum to a string here
	Nationality          string      `json:"nationality"`
	StudyDegree          string      `json:"studyDegree"` // for simplicity we map the enum to a string here
	StudyProgram         string      `json:"studyProgram"`
	CurrentSemester      pgtype.Int4 `json:"currentSemester"`
}

type GetAllCPPsForCoursePhase struct {
	CoursePhaseID         uuid.UUID `json:"coursePhaseID"`
	PassStatus            string    `json:"passStatus"`
	CourseParticipationID uuid.UUID `json:"courseParticipationID"`
	RestrictedData        MetaData  `json:"restrictedData"`
	StudentReadableData   MetaData  `json:"studentReadableData"`
	PrevData              MetaData  `json:"prevData"`
	Student               Student   `json:"student"`
}

type CoursePhaseParticipationsWithResolutions struct {
	Participations []GetAllCPPsForCoursePhase `json:"participations"`
	Resolutions    []Resolution               `json:"resolutions"`
}

type PrevCoursePhaseData struct {
	PrevData    MetaData     `json:"prevData"`
	Resolutions []Resolution `json:"resolutions"`
}

// buildURL constructs the request URL for a given resolution.
// extraPaths (such as a courseParticipationID) can be appended.
func buildURL(resolution Resolution, extraPaths ...string) string {
	base := fmt.Sprintf("%s/course_phase/%s/%s", resolution.BaseURL, resolution.CoursePhaseID, getEndpointPath(resolution.EndpointPath))
	if len(extraPaths) > 0 {
		base = fmt.Sprintf("%s/%s", base, strings.Join(extraPaths, "/"))
	}
	return base
}

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
func ResolveParticipation(authHeader string, resolution Resolution, courseParticipationID uuid.UUID) (interface{}, error) {
	url := buildURL(resolution, courseParticipationID.String())
	data, err := FetchJSON(url, authHeader)
	if err != nil {
		return nil, err
	}

	return parseAndValidate(data, resolution.DtoName)
}

// ResolveCoursePhaseData resolves data for a course phase.
func ResolveCoursePhaseData(authHeader string, resolution Resolution) (interface{}, error) {
	url := buildURL(resolution)
	data, err := FetchJSON(url, authHeader)
	if err != nil {
		return nil, err
	}

	return parseAndValidate(data, resolution.DtoName)
}

// ResolveAllParticipations resolves data for all participations and returns a map keyed by courseParticipationID.
func ResolveAllParticipations(authHeader string, resolution Resolution) (map[uuid.UUID]interface{}, error) {
	url := buildURL(resolution)
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
func FetchAndMergeParticipationsWithResolutions(coreURL string, authHeader string, coursePhaseID uuid.UUID) ([]GetAllCPPsForCoursePhase, error) {
	url := fmt.Sprintf("%s/api/course_phases/%s/participations", coreURL, coursePhaseID)
	data, err := FetchJSON(url, authHeader)
	if err != nil {
		return nil, err
	}

	var cppWithRes CoursePhaseParticipationsWithResolutions
	if err := json.Unmarshal(data, &cppWithRes); err != nil {
		return nil, err
	}

	for _, res := range cppWithRes.Resolutions {
		resolvedData, err := ResolveAllParticipations(authHeader, res)
		if err != nil {
			return nil, err
		}

		for idx, participation := range cppWithRes.Participations {
			if data, exists := resolvedData[participation.CourseParticipationID]; exists {
				if participation.PrevData == nil {
					participation.PrevData = make(MetaData)
				}
				participation.PrevData[res.DtoName] = data
				cppWithRes.Participations[idx] = participation
			}
		}
	}

	return cppWithRes.Participations, nil
}

func FetchAndMergeCoursePhaseWithResolution(coreURL string, authHeader string, coursePhaseID uuid.UUID) (MetaData, error) {
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
		cpWithRes.PrevData = make(MetaData)
	}

	for _, res := range cpWithRes.Resolutions {
		resolvedData, err := ResolveCoursePhaseData(authHeader, res)
		if err != nil {
			return nil, err
		}
		cpWithRes.PrevData[res.DtoName] = resolvedData
	}
	return cpWithRes.PrevData, nil
}

// getEndpointPath trims leading and trailing slashes from the endpoint path.
func getEndpointPath(endpointPath string) string {
	return strings.Trim(endpointPath, "/")
}
