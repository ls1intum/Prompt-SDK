package promptTypes

import "github.com/google/uuid"

type Person struct {
	ID        uuid.UUID `json:"id" binding:"uuid"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}
