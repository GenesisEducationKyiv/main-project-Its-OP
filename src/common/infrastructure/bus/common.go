package bus

import (
	"github.com/google/uuid"
)

type Command struct {
	ID   uuid.UUID `json:"id"`
	Body string    `json:"body"`
}
