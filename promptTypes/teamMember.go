package promptTypes

import "github.com/google/uuid"

type TeamMember struct {
	CourseParticipationID uuid.UUID `json:"courseParticipationID"`
	StudentName           string    `json:"studentName"`
}
