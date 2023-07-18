package infrastructure

import (
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

const (
	LogLevelInfo  = "Info"
	LogLevelDebug = "Debug"
	LogLevelError = "Error"
)

type Logger struct {
	commandBus *cqrs.CommandBus
}

func NewLogger(commandBus *cqrs.CommandBus) *Logger {
	return &Logger{commandBus: commandBus}
}

func (l *Logger) LogInformation(message string) error {
	logCommand := commands.NewLogCommand(fmt.Sprintf("INFO - %s", message), LogLevelInfo)
	return l.commandBus.Send(context.Background(), logCommand)
}

func (l *Logger) LogDebug(message string) error {
	logCommand := commands.NewLogCommand(fmt.Sprintf("DEBUG - %s", message), LogLevelDebug)
	return l.commandBus.Send(context.Background(), logCommand)
}

func (l *Logger) LogError(err error, message string) error {
	logCommand := commands.NewLogCommand(fmt.Sprintf("ERROR - %s; {errMsg: %s}", message, err.Error()), LogLevelError)
	return l.commandBus.Send(context.Background(), logCommand)
}
