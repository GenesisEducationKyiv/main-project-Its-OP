package logger

import (
	"btcRate/common/application"
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"fmt"
	"golang.org/x/exp/slog"
)

type AsyncLogger struct {
	commandBus          application.ICommandBus
	logCommandValidator application.IValidator[commands.LogCommand]
}

func NewAsyncLogger(commandBus application.ICommandBus, logCommandValidator application.IValidator[commands.LogCommand]) *AsyncLogger {
	return &AsyncLogger{commandBus: commandBus, logCommandValidator: logCommandValidator}
}

func (l *AsyncLogger) Info(message string, args ...any) {
	l.send(commands.NewLogCommand(message, args, slog.LevelInfo))
}

func (l *AsyncLogger) Debug(message string, args ...any) {
	l.send(commands.NewLogCommand(message, args, slog.LevelDebug))
}

func (l *AsyncLogger) Error(message string, args ...any) {
	l.send(commands.NewLogCommand(message, args, slog.LevelError))
}

func (l *AsyncLogger) send(c *commands.LogCommand) {
	if err := l.logCommandValidator.Validate(c); err != nil {
		fmt.Println("log command is invalid:", err.Error())
	}

	if err := l.commandBus.Send(context.Background(), c); err != nil {
		fmt.Println("failed to send log command:", err.Error())
	}
}
