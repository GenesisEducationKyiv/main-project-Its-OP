package bus

import (
	"github.com/google/uuid"
)

type Command struct {
	ID   uuid.UUID `json:"id"`
	Body string    `json:"body"`
}

func NewCommand(id uuid.UUID, body string) Command {
	return Command{ID: id, Body: body}
}

type ICommandHandler interface {
	GetName() string
	Handle(cmd Command)
}
