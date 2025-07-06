package promptTypes

import "github.com/google/uuid"

type CopyRequest struct {
	SourceCoursePhaseID uuid.UUID `json:"sourceCoursePhaseID"`
	TargetCoursePhaseID uuid.UUID `json:"targetCoursePhaseID"`
}
