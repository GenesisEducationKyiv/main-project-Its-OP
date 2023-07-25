package command_handlers

import (
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"fmt"
	"os"
)

type ErrorCommandHandler struct {
}

func (h ErrorCommandHandler) HandlerName() string {
	return ErrorLogCommandHandlerName
}

func (h ErrorCommandHandler) NewCommand() interface{} {
	return &commands.LogCommand{}
}

func (h ErrorCommandHandler) Handle(_ context.Context, cmd interface{}) error {
	logCommand := cmd.(*commands.LogCommand)
	_, err := fmt.Fprint(os.Stderr, fmt.Sprintf("%s\n", logCommand.LogData))

	return err
}
