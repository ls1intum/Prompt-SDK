package promptTypes

import "github.com/google/uuid"

// Person represents an individual person in the Prompt system.
// This type is used to identify and store basic information about users,
// students, lecturers, and other participants within the platform.
type Person struct {
	// ID is the unique identifier for the person using UUID v4 format.
	// This field is required and must be a valid UUID.
	ID uuid.UUID `json:"id" binding:"uuid"`

	// FirstName stores the first name as provided during initial application or profile updates.
	FirstName string `json:"firstName"`

	// LastName stores the last name as provided during initial application or profile updates.
	LastName string `json:"lastName"`
}
