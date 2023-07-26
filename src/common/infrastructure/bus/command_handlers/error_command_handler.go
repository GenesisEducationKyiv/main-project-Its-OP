package command_handlers

import (
	"btcRate/common/application"
	"btcRate/common/infrastructure/bus/commands"
	"context"
)

type ErrorCommandHandler struct {
	logger application.ILogger
}

func (h ErrorCommandHandler) HandlerName() string {
	return ErrorLogCommandHandlerName
}

func (h ErrorCommandHandler) NewCommand() interface{} {
	return &commands.LogCommand{}
}

func (h ErrorCommandHandler) Handle(_ context.Context, cmd interface{}) error {
	logCommand := cmd.(*commands.LogCommand)
	return h.logger.Error(logCommand.LogMessage, logCommand.LogAttributes)
}
