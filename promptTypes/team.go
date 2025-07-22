package promptTypes

import "github.com/google/uuid"

type Team struct {
	ID      uuid.UUID `json:"id" binding:"uuid"`
	Name    string    `json:"name" binding:"required"`
	Members []Person  `json:"members" binding:"dive"`
	Tutors  []Person  `json:"tutors" binding:"dive"`
}
