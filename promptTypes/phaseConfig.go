package promptTypes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PhaseConfigHandler defines the interface that any module must implement
// to handle a course phase configuration request initiated by the core system.
type PhaseConfigHandler interface {
	// HandlePhaseConfig processes the retrieval of configuration settings
	// for the specified course phase.
	// It returns a map of configuration settings and their statuses.
	// The map keys are configuration names, and the values are booleans indicating whether the
	// setting is configured or is missing in the phase.
	HandlePhaseConfig() (map[string]bool, error)
}

// RegisterConfigEndpoint registers the standardized GET /config endpoint on the given router group.
// It applies the provided authorization middleware and delegates handling to the provided PhaseConfigHandler.
// Example endpoint path:
//
//	GET /self-team-allocation/api/course_phase/:coursePhaseID/config
func RegisterConfigEndpoint(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, handler PhaseConfigHandler) {
	router.GET("/config", authMiddleware, func(c *gin.Context) {
		response, err := handler.HandlePhaseConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})
}
