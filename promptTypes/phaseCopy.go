package promptTypes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PhaseCopyRequest represents the payload used to trigger the copying of a course phase.
// This request initiates the duplication of course phase data, settings, and state
// from a source phase to a target phase, enabling course reuse and template functionality.
type PhaseCopyRequest struct {
	// SourceCoursePhaseID is the unique identifier of the course phase to copy from.
	// This phase serves as the template and its data will be duplicated.
	SourceCoursePhaseID uuid.UUID `json:"sourceCoursePhaseID"`

	// TargetCoursePhaseID is the unique identifier of the course phase to copy into.
	// This phase will receive the copied data and configurations.
	TargetCoursePhaseID uuid.UUID `json:"targetCoursePhaseID"`
}

// PhaseCopyHandler defines the interface that modules must implement to support phase copying.
// Any module that stores course phase-specific data should implement this interface
// to ensure their data is properly copied when course phases are duplicated.
type PhaseCopyHandler interface {
	// HandlePhaseCopy processes the copying of module-specific data from the source
	// course phase to the target course phase. Implementations should:
	//   - Copy relevant settings and configurations
	//   - Duplicate user data where appropriate
	//   - Update references and relationships
	//   - Maintain data integrity and consistency
	//
	// Returns an error if the copy operation fails for any reason.
	HandlePhaseCopy(c *gin.Context, req PhaseCopyRequest) error
}

// RegisterCopyEndpoint registers the standardized POST /copy endpoint for phase copying.
// This function sets up the HTTP endpoint that modules can use to expose phase copy functionality
// in a consistent manner across the Prompt platform.
//
// The endpoint handles:
//   - JSON request parsing and validation
//   - Authentication through the provided middleware
//   - Error handling and standardized responses
//   - Success confirmation messages
//
// Example endpoint path: POST /self-team-allocation/api/course_phase/:coursePhaseID/copy
//
// Parameters:
//   - router: The Gin router group where the endpoint will be registered
//   - authMiddleware: Authentication middleware to protect the endpoint
//   - handler: Implementation of PhaseCopyHandler that performs the actual copying
func RegisterCopyEndpoint(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, handler PhaseCopyHandler) {
	router.POST("/copy", authMiddleware, func(c *gin.Context) {
		var req PhaseCopyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := handler.HandlePhaseCopy(c, req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "Course phase copied successfully"})
	})
}
