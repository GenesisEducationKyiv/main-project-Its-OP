package commands

import (
	"github.com/google/uuid"
)

type LogCommand struct {
	ID  uuid.UUID
	Log string
}

func NewLogCommand(id uuid.UUID, body string) *LogCommand {
	return &LogCommand{ID: id, Log: body}
}
