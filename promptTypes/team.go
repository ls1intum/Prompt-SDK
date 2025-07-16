package promptTypes

import "github.com/google/uuid"

type Team struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Members []Person  `json:"members"`
	Tutors  []Person  `json:"tutors"`
}
