package keycloakTokenVerifier

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ls1intum/prompt-sdk/keycloakTokenVerifier/keycloakCoreRequests"
	log "github.com/sirupsen/logrus"
)

// Important: This requires a CoursePhaseID as a parameter.
func getLecturerAndEditorRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		coursePhaseID, err := uuid.Parse(c.Param("coursePhaseID"))
		if err != nil {
			log.Error("Error parsing coursePhaseID:", err)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if coursePhaseID == uuid.Nil {
			log.Error("Invalid coursePhaseID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("coursePhaseID missing"))
			return
		}

		// TODO: Wrap this around a caching component
		// retrieve the relevant roles from the core
		tokenMapping, err := keycloakCoreRequests.SendCoursePhaseRoleMappingRequest(KeycloakTokenVerifierSingleton.CoreURL, c.GetHeader("Authorization"), coursePhaseID)
		if err != nil {
			log.Error("Error getting course roles:", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// get roles from the context
		token_student, ok := GetTokenStudent(c)
		if !ok {
			log.Error("Error getting token student")
			_ = c.AbortWithError(http.StatusInternalServerError, ErrStudentNotInContext)
			return
		}
		userRoles := token_student.Roles

		// filter out the roles relevant for the current course phase
		isLecturer := userRoles[tokenMapping.CourseLecturerRole]
		isEditor := userRoles[tokenMapping.CourseEditorRole]

		// DEPRECATED: Keep this for backwards compatibility
		c.Set("isLecturer", isLecturer)
		c.Set("isEditor", isEditor)
		c.Set("customRolePrefix", tokenMapping.CustomRolePrefix)

		token_student.IsLecturer = isLecturer
		token_student.IsEditor = isEditor
		token_student.CustomRolePrefix = tokenMapping.CustomRolePrefix
		SetTokenStudent(c, token_student)
	}
}
