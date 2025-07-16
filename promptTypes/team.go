package promptTypes

import "github.com/google/uuid"

type Team struct {
	ID      uuid.UUID    `json:"id"`
	Name    string       `json:"name"`
	Members []TeamMember `json:"members"`
	Tutors  []TeamMember `json:"tutors"`
}
