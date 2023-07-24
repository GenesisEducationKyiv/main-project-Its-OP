package infrastructure

import (
	"btcRate/common/application"
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"fmt"
)

const (
	LogLevelInfo  = "Info"
	LogLevelDebug = "Debug"
	LogLevelError = "Error"
)

type Logger struct {
	commandBus application.ICommandBus
}

func NewLogger(commandBus application.ICommandBus) *Logger {
	return &Logger{commandBus: commandBus}
}

func (l *Logger) Info(message string) error {
	logCommand := commands.NewLogCommand(fmt.Sprintf("INFO - %s", message), LogLevelInfo)
	return l.commandBus.Send(context.Background(), logCommand)
}

func (l *Logger) Debug(message string) error {
	logCommand := commands.NewLogCommand(fmt.Sprintf("DEBUG - %s", message), LogLevelDebug)
	return l.commandBus.Send(context.Background(), logCommand)
}

func (l *Logger) Error(err error, message string) error {
	logCommand := commands.NewLogCommand(fmt.Sprintf("ERROR - %s; {errMsg: %s}", message, err.Error()), LogLevelError)
	return l.commandBus.Send(context.Background(), logCommand)
}
