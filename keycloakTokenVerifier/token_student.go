package keycloakTokenVerifier

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const tokenStudentContextKey = "token_student"

type TokenStudent struct {
	Roles               map[string]bool
	ID                  string
	Email               string
	MatriculationNumber string
	UniversityLogin     string
	FirstName           string
	LastName            string

	IsStudentOfCourse      bool
	IsStudentOfCoursePhase bool
	CourseParticipationID  uuid.UUID

	IsLecturer       bool
	IsEditor         bool
	CustomRolePrefix string
}

func GetTokenStudent(c *gin.Context) (TokenStudent, bool) {
	if tokenStudent, exists := c.Get(tokenStudentContextKey); exists {
		return tokenStudent.(TokenStudent), true
	}
	return TokenStudent{}, false
}

func SetTokenStudent(c *gin.Context, tokenStudent TokenStudent) {
	c.Set(tokenStudentContextKey, tokenStudent)
}
