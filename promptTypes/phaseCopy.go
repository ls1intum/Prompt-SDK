package promptTypes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PhaseCopyRequest struct {
	SourceCoursePhaseID uuid.UUID `json:"sourceCoursePhaseID"`
	TargetCoursePhaseID uuid.UUID `json:"targetCoursePhaseID"`
}

type PhaseCopyHandler interface {
	HandlePhaseCopy(ctx context.Context, req PhaseCopyRequest) error
}

func RegisterCopyEndpoint(router *gin.RouterGroup, handler PhaseCopyHandler) {
	router.POST("/copy", func(c *gin.Context) {
		var req PhaseCopyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := handler.HandlePhaseCopy(c.Request.Context(), req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})
}
