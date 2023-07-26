package logger

import (
	"btcRate/common/application"
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"golang.org/x/exp/slog"
)

type AsyncLogger struct {
	commandBus          application.ICommandBus
	logCommandValidator application.IValidator[commands.LogCommand]
}

func NewLogger(commandBus application.ICommandBus) *AsyncLogger {
	return &AsyncLogger{commandBus: commandBus}
}

func (l *AsyncLogger) Info(message string, args ...any) error {
	return l.commandBus.Send(context.Background(), commands.NewLogCommand(message, args, slog.LevelInfo))
}

func (l *AsyncLogger) Debug(message string, args ...any) error {
	return l.commandBus.Send(context.Background(), commands.NewLogCommand(message, args, slog.LevelDebug))
}

func (l *AsyncLogger) Error(message string, args ...any) error {
	return l.commandBus.Send(context.Background(), commands.NewLogCommand(message, args, slog.LevelError))
}

func (l *AsyncLogger) send(c *commands.LogCommand) error {
	if err := l.logCommandValidator.Validate(c); err != nil {
		return err
	}

	return l.commandBus.Send(context.Background(), c)
}
