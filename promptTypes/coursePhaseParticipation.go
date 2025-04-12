package promptTypes

import "github.com/google/uuid"

type CoursePhaseParticipationWithStudent struct {
	CoursePhaseID         uuid.UUID `json:"coursePhaseID"`
	PassStatus            string    `json:"passStatus"`
	CourseParticipationID uuid.UUID `json:"courseParticipationID"`
	RestrictedData        MetaData  `json:"restrictedData"`
	StudentReadableData   MetaData  `json:"studentReadableData"`
	PrevData              MetaData  `json:"prevData"`
	Student               Student   `json:"student"`
}
