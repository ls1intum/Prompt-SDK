package promptSDK

import "github.com/google/uuid"

// TODO

type Resolution struct {
	DtoName       string
	BaseURL       string
	EndpointPath  string
	CoursePhaseID uuid.UUID
}

func ResolveParticipation(courseParticipationID uuid.UUID, resolution Resolution) (interface{}, error) {
	// TODO
	return nil, nil
}

func ResolveAllParticipations(resolution Resolution) (interface{}, error) {
	// TODO
	return nil, nil
}

func ResolveCoursePhaseData(resolution Resolution) (interface{}, error) {
	// TODO
	return nil, nil
}
