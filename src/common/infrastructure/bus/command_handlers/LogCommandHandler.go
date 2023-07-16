package command_handlers

import (
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"fmt"
)

type LogCommandHandler struct {
}

func (h LogCommandHandler) HandlerName() string {
	return "LogCommandHandler"
}

func (h LogCommandHandler) NewCommand() interface{} {
	return &commands.LogCommand{}
}

func (h LogCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	logCommand := cmd.(*commands.LogCommand)
	fmt.Print(logCommand.Log)

	return nil
}
