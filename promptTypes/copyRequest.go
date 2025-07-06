package promptTypes

import "github.com/google/uuid"

type CopyRequest struct {
	CoursePhaseIDOld uuid.UUID `json:"coursePhaseIDOld"`
	CoursePhaseIDNew uuid.UUID `json:"coursePhaseIDNew"`
}
