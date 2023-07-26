package command_handlers

import (
	"btcRate/common/application"
	"btcRate/common/infrastructure/bus/commands"
	"context"
)

type ErrorCommandHandler struct {
	logger application.ILogger
}

func NewErrorCommandHandler(logger application.ILogger) *ErrorCommandHandler {
	return &ErrorCommandHandler{logger: logger}
}

func (h ErrorCommandHandler) HandlerName() string {
	return ErrorLogCommandHandlerName
}

func (h ErrorCommandHandler) NewCommand() interface{} {
	return &commands.LogCommand{}
}

func (h ErrorCommandHandler) Handle(_ context.Context, cmd interface{}) error {
	logCommand := cmd.(*commands.LogCommand)
	h.logger.Error(logCommand.LogMessage, logCommand.LogAttributes)
	return nil
}
