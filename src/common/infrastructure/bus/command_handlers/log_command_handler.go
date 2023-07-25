package command_handlers

import (
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"fmt"
)

type LogCommandHandler struct {
}

func (h LogCommandHandler) HandlerName() string {
	return LogCommandHandlerName
}

func (h LogCommandHandler) NewCommand() interface{} {
	return &commands.LogCommand{}
}

func (h LogCommandHandler) Handle(_ context.Context, cmd interface{}) error {
	logCommand := cmd.(*commands.LogCommand)
	fmt.Printf("%s\n", logCommand.LogData)

	return nil
}
