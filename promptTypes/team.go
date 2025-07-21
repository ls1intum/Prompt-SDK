package promptTypes

import "github.com/google/uuid"

// Team represents a group of students and tutors working together in a course phase.
// This type is used for team-based activities, group projects, and collaborative learning
// where multiple participants need to be organized and managed as a cohesive unit.
type Team struct {
	ID uuid.UUID `json:"id" binding:"uuid"`

	// Name is the display name of the team, used for identification and communication.
	// This could be assigned automatically (e.g., "Team 1") or chosen by team members.
	Name string `json:"name" binding:"required"`

	// Members contains the list of students who are part of this team.
	// These are the primary participants who will be working together on team activities.
	Members []Person `json:"members" binding:"dive"`

	// Tutors contains the list of tutors, coaches, or teaching assistants assigned to guide this team.
	// Tutors provide mentorship, feedback, and academic support to the team members.
	Tutors []Person `json:"tutors" binding:"dive"`
}
