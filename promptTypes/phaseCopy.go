package promptTypes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PhaseCopyRequest struct {
	SourceCoursePhaseID uuid.UUID `json:"sourceCoursePhaseID"`
	TargetCoursePhaseID uuid.UUID `json:"targetCoursePhaseID"`
}

type PhaseCopyHandler interface {
	HandlePhaseCopy(c *gin.Context, req PhaseCopyRequest) error
}

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
