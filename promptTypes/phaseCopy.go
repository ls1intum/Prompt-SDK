package promptTypes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PhaseCopyRequest represents the payload used to trigger
// the copying of a course phase from a source to a target phase.
type PhaseCopyRequest struct {
	SourceCoursePhaseID uuid.UUID `json:"sourceCoursePhaseID"` // ID of the phase to copy from
	TargetCoursePhaseID uuid.UUID `json:"targetCoursePhaseID"` // ID of the phase to copy into
}

// PhaseCopyHandler defines the interface that any module must implement
// to handle a phase copy operation initiated by the core system.
type PhaseCopyHandler interface {
	// HandlePhaseCopy processes the copying of internal state/data
	// from the source course phase to the target course phase.
	HandlePhaseCopy(c *gin.Context, req PhaseCopyRequest) error
}

// RegisterCopyEndpoint registers the standardized POST /copy endpoint on the given router group.
// It applies the provided authorization middleware and delegates handling to the provided PhaseCopyHandler.
//
// Example endpoint path:
//
//	POST /self-team-allocation/api/course_phase/:coursePhaseID/copy
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
