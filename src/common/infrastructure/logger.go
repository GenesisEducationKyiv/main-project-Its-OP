package infrastructure

import (
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Logger struct {
	commandBus *cqrs.CommandBus
}

func NewLogger(commandBus *cqrs.CommandBus) *Logger {
	return &Logger{commandBus: commandBus}
}

func (l *Logger) LogInformation(message string) error {
	logCommand := &commands.LogCommand{Log: fmt.Sprintf("INFO - %s", message)}
	return l.commandBus.Send(context.Background(), logCommand)
}

func (l *Logger) LogDebug(message string) error {
	logCommand := &commands.LogCommand{Log: fmt.Sprintf("DEBUG - %s", message)}
	return l.commandBus.Send(context.Background(), logCommand)
}

func (l *Logger) LogError(err error, message string) error {
	logCommand := &commands.LogCommand{Log: fmt.Sprintf("ERROR - %s; {errMsg: %s}", message, err.Error())}
	return l.commandBus.Send(context.Background(), logCommand)
}
