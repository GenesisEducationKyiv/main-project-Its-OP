package command_handlers

import (
	"btcRate/common/application"
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"golang.org/x/exp/slog"
)

type LogCommandHandler struct {
	logger application.ILogger
}

func (h LogCommandHandler) HandlerName() string {
	return LogCommandHandlerName
}

func (h LogCommandHandler) NewCommand() interface{} {
	return &commands.LogCommand{}
}

func (h LogCommandHandler) Handle(_ context.Context, cmd interface{}) error {
	logCommand := cmd.(*commands.LogCommand)
	switch logCommand.LogLevel {
	case slog.LevelInfo:
		return h.logger.Info(logCommand.LogMessage, logCommand.LogAttributes)

	case slog.LevelDebug:
		return h.logger.Debug(logCommand.LogMessage, logCommand.LogAttributes)

	case slog.LevelError:
		return h.logger.Error(logCommand.LogMessage, logCommand.LogAttributes)

	default:
		return h.logger.Error("cannot handle log level", "log_level", logCommand.LogLevel, "command")
	}
}
