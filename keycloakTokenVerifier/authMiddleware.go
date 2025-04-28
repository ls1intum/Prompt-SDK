package keycloakTokenVerifier

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// AuthenticationMiddleware creates a composite middleware which always
// applies KeycloakMiddleware first and conditionally chains additional
// middlewares based on the allowed roles:
//   - If allowedRoles contains "Lecturer", "Editor" or a custom role name
//     (any value other than "Admin" or "Student"), then it calls GetLecturerAndEditorRole.
//     For custom roles the middleware checks if the user's roles include customRolePrefix+customRole.
//   - If allowedRoles contains "Student", then it calls IsStudentOfCoursePhaseMiddleware.
func AuthenticationMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Always run Keycloak middleware first.
		KeycloakMiddleware()(c)
		if c.IsAborted() {
			return
		}

		allowedSet := buildAllowedRolesSet(allowedRoles)

		tokenUser, ok := GetTokenUser(c)
		if !ok {
			log.Error("Error getting token student")
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrUserNotInContext)
			return
		}
		userRoles := tokenUser.Roles

		// 1.) Directly grant access for PROMPT_Admin or PROMPT_Lecturer.
		if checkDirectRole(PromptAdmin, allowedSet, userRoles) ||
			checkDirectRole(PromptLecturer, allowedSet, userRoles) {
			c.Next()
			return
		}

		// This allows to use the middleware without coursePhaseID, if only PROMPT_Admin & PROMPT_Lecturer are allowed.
		if onlyContainsAdminAndLecturer(allowedSet) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not authenticate"})
			return
		}

		// 2.) Check for Lecturer, Editor, or custom group roles.
		if requiresLecturerOrCustom(allowedSet, allowedRoles) {
			getLecturerAndEditorRole()(c)
			if c.IsAborted() {
				return
			}

			tokenUser, ok = GetTokenUser(c)
			if !ok {
				log.Error("Error refreshing the token student")
				c.AbortWithStatusJSON(http.StatusUnauthorized, ErrUserNotInContext)
				return
			}

			if _, allowed := allowedSet[CourseLecturer]; allowed && tokenUser.IsLecturer {
				c.Next()
				return
			}

			if _, allowed := allowedSet[CourseEditor]; allowed && tokenUser.IsEditor {
				c.Next()
				return
			}

			if containsCustomRoleName(allowedRoles...) {
				prefix := tokenUser.CustomRolePrefix

				for _, role := range allowedRoles {
					if userRoles[prefix+role] {
						c.Next()
						return
					}
				}
			}
		}

		// 3.) Check for Student.
		if _, allowed := allowedSet[CourseStudent]; allowed {
			isStudentOfCoursePhaseMiddleware()(c)
			if c.IsAborted() {
				return
			}

			tokenUser, ok = GetTokenUser(c)
			if !ok {
				log.Error("Error refreshing the token student")
				c.AbortWithStatusJSON(http.StatusUnauthorized, ErrUserNotInContext)
				return
			}

			if tokenUser.IsStudentOfCourse {
				c.Next()
				return
			}
		}

		// Access denied.
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not authenticate"})
	}
}

// buildAllowedRolesSet creates a lookup set from a slice of roles.
func buildAllowedRolesSet(roles []string) map[string]struct{} {
	set := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		set[role] = struct{}{}
	}
	return set
}

// checkDirectRole returns true if a specific role is both allowed and present in the user roles.
func checkDirectRole(role string, allowedSet map[string]struct{}, userRoles map[string]bool) bool {
	if _, allowed := allowedSet[role]; allowed && userRoles[role] {
		return true
	}
	return false
}

// requiresLecturerOrCustom determines if additional checks for lecturer, editor,
// or custom roles are needed based on the allowed roles.
func requiresLecturerOrCustom(allowedSet map[string]struct{}, roles []string) bool {
	_, hasLecturer := allowedSet[CourseLecturer]
	_, hasEditor := allowedSet[CourseEditor]
	return hasLecturer || hasEditor || containsCustomRoleName(roles...)
}

func containsCustomRoleName(allowedRoles ...string) bool {
	nonCustomRoles := []string{PromptAdmin, PromptLecturer, CourseLecturer, CourseEditor, CourseStudent}

	for _, role := range allowedRoles {
		if !slices.Contains(nonCustomRoles, role) {
			return true
		}
	}

	return false
}

// onlyContainsAdminAndLecturer returns true if the allowedSet only contains
// "PROMPT_Admin" and/or "PROMPT_Lecturer".
func onlyContainsAdminAndLecturer(allowedSet map[string]struct{}) bool {
	for role := range allowedSet {
		if role != PromptAdmin && role != PromptLecturer {
			return false
		}
	}
	return true
}
