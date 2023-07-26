package decorators

import (
	"btcRate/common/application"
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"time"
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
	processingId := uuid.New()
	h.logger.Info("command processing started", "commandName", commandName, "processingId", processingId)

	start := time.Now()

	err := h.handler.Handle(context, cmd)

	elapsed := fmt.Sprintf("%dms", time.Since(start).Milliseconds())

	if err == nil {
		h.logger.Info("command processing finished", "status", "Success", "processingId", processingId, "processingTime", elapsed)
	} else {
		h.logger.Error("command processing finished", "status", "Failure", "processingId", processingId, "processingTime", elapsed, "error", err.Error())
	}

	return err
}
