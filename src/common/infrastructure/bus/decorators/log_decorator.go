package decorators

import (
	"btcRate/common/application"
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type LogDecorator struct {
	handler      cqrs.CommandHandler
	generateName func(v interface{}) string
	logger       application.ILogger
}

func NewLoggedCommandHandler(handler cqrs.CommandHandler, generateName func(v interface{}) string, logger application.ILogger) LogDecorator {
	return LogDecorator{handler: handler, generateName: generateName, logger: logger}
}

func (h LogDecorator) HandlerName() string {
	return h.handler.HandlerName()
}

func (h LogDecorator) NewCommand() interface{} {
	return h.handler.NewCommand()
}

func (h LogDecorator) Handle(context context.Context, cmd interface{}) error {
	commandName := h.generateName(cmd)
	logErr := h.logger.Info("command processing started", "commandName", commandName)
	if logErr != nil {
		fmt.Println("failed to log", logErr.Error())
	}

	err := h.handler.Handle(context, cmd)

	if err == nil {
		logErr = h.logger.Info("command processing finished", "status", "Success")
	} else {
		logErr = h.logger.Error("command processing finished", "status", "Failure", "error", err.Error())
	}

	if logErr != nil {
		fmt.Println("failed to log", err.Error())
	}

	return err
}
