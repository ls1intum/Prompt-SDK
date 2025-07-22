package promptTypes

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Student struct {
	Person
	Email                string      `json:"email" binding:"email"`
	MatriculationNumber  string      `json:"matriculationNumber" binding:"matriculationNumber"`
	UniversityLogin      string      `json:"universityLogin" binding:"universityLogin"`
	HasUniversityAccount bool        `json:"hasUniversityAccount"`
	Gender               Gender      `json:"gender" binding:"oneof=male female diverse prefer_not_to_say"`
	Nationality          string      `json:"nationality"`
	StudyDegree          StudyDegree `json:"studyDegree" binding:"oneof=bachelor master"`
	StudyProgram         string      `json:"studyProgram"`
	CurrentSemester      pgtype.Int4 `json:"currentSemester"`
}
