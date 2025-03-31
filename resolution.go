package promptSDK

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func FetchAndMergeParticipationsWithResolutions(coreURL string, authHeader string, coursePhaseID uuid.UUID) ([]GetAllCPPsForCoursePhase, error) {
	url := fmt.Sprintf("%s/course_phases/%s/participations", coreURL, coursePhaseID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cppWithRes CoursePhaseParticipationsWithResolutions
	if err := json.Unmarshal(data, &cppWithRes); err != nil {
		return nil, err
	}

	for _, res := range cppWithRes.Resolutions {
		resolvedDataMap, err := ResolveAllParticipations(authHeader, res)
		if err != nil {
			return nil, err
		}

		resolvedData, ok := resolvedDataMap.(map[uuid.UUID]interface{})
		if !ok {
			log.Error("Failed to cast resolved data to map[uuid.UUID]interface{}")
			continue
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

func ResolveParticipation(authHeader string, resolution Resolution, courseParticipationID uuid.UUID) (interface{}, error) {
	url := fmt.Sprintf("%s/coursePhase/%s/%s/%s", resolution.BaseURL, resolution.CoursePhaseID, resolution.EndpointPath, courseParticipationID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result[resolution.DtoName], nil
}

func ResolveAllParticipations(authHeader string, resolution Resolution) (interface{}, error) {
	url := fmt.Sprintf("%s/coursePhase/%s/%s", resolution.BaseURL, resolution.CoursePhaseID, resolution.EndpointPath)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	formattedResults := make(map[uuid.UUID]interface{})
	for _, item := range results {
		idStr, ok := item["coursePhaseParticipationID"].(string)
		if !ok {
			continue
		}
		participationID, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		formattedResults[participationID] = item[resolution.DtoName]
	}

	return formattedResults, nil
}

func ResolveCoursePhaseData(authHeader string, resolution Resolution) (interface{}, error) {
	url := fmt.Sprintf("%s/coursePhase/%s/%s", resolution.BaseURL, resolution.CoursePhaseID, resolution.EndpointPath)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result[resolution.DtoName], nil
}
