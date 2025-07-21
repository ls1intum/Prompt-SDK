package promptTypes

import "github.com/google/uuid"

// CoursePhaseParticipationWithStudent represents a student's participation in a specific course phase.
// This type combines participation metadata with the student's information, providing a complete
// view of how a student is engaged with a particular phase of a course.
type CoursePhaseParticipationWithStudent struct {
	// CoursePhaseID is the unique identifier of the course phase the student is participating in.
	CoursePhaseID uuid.UUID `json:"coursePhaseID"`

	// PassStatus indicates the student's current status in this course phase.
	// Common values are usually "passed", "failed", and "not_assessed".
	PassStatus string `json:"passStatus"`

	// CourseParticipationID links this phase participation to the overall course participation.
	CourseParticipationID uuid.UUID `json:"courseParticipationID"`

	// RestrictedData contains sensitive metadata that should only be accessible to authorized users
	// such as instructors, admins, or the system itself.
	RestrictedData MetaData `json:"restrictedData"`

	// StudentReadableData contains metadata that can be safely shared with the student,
	// such as feedback, scores, or progress information.
	StudentReadableData MetaData `json:"studentReadableData"`

	// PrevData contains metadata from previous phases, passed via course-phase communication
	PrevData MetaData `json:"prevData"`

	// Student contains the complete student information associated with this participation.
	Student Student `json:"student"`
}
