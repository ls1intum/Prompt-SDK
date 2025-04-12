package promptTypes

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Student struct {
	ID                   uuid.UUID   `json:"id"`
	FirstName            string      `json:"firstName"`
	LastName             string      `json:"lastName"`
	Email                string      `json:"email"`
	MatriculationNumber  string      `json:"matriculationNumber"`
	UniversityLogin      string      `json:"universityLogin"`
	HasUniversityAccount bool        `json:"hasUniversityAccount"`
	Gender               Gender      `json:"gender"`
	Nationality          string      `json:"nationality"`
	StudyDegree          StudyDegree `json:"studyDegree"`
	StudyProgram         string      `json:"studyProgram"`
	CurrentSemester      pgtype.Int4 `json:"currentSemester"`
}
