package commands

import "golang.org/x/exp/slog"

type LogCommand struct {
	LogMessage    string
	LogAttributes []any
	LogLevel      slog.Level
}

func NewLogCommand(message string, attributes []any, level slog.Level) *LogCommand {
	return &LogCommand{LogMessage: message, LogAttributes: attributes, LogLevel: level}
}

func (c *LogCommand) Reset() {
	c.LogMessage = ""
	c.LogAttributes = make([]any, 16)
}

func (c *LogCommand) String() string {
	return c.LogMessage
}

func (c *LogCommand) ProtoMessage() {}
