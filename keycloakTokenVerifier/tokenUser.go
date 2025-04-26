package keycloakTokenVerifier

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const tokenUserContextKey = "tokenUser"

var ErrStudentNotInContext = errors.New("user not found in context")

// TokenUser encapsulates a user's authentication information, including roles,
// personal identifiers, and permissions within the system.
type TokenUser struct {
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

func GetTokenUser(c *gin.Context) (TokenUser, bool) {
	if tokenUser, exists := c.Get(tokenUserContextKey); exists {
		ts, ok := tokenUser.(TokenUser)
		if !ok {
			return TokenUser{}, false
		}
		return ts, true
	}
	return TokenUser{}, false
}

func SetTokenUser(c *gin.Context, tokenUser TokenUser) {
	c.Set(tokenUserContextKey, tokenUser)
}
