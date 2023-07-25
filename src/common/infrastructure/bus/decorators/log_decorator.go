package decorators

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"time"
)

type LogDecorator struct {
	handler      cqrs.CommandHandler
	generateName func(v interface{}) string
}

func NewLoggedCommandHandler(handler cqrs.CommandHandler, generateName func(v interface{}) string) LogDecorator {
	return LogDecorator{handler: handler, generateName: generateName}
}

func (h LogDecorator) HandlerName() string {
	return h.handler.HandlerName()
}

func (h LogDecorator) NewCommand() interface{} {
	return h.handler.NewCommand()
}

func (h LogDecorator) Handle(context context.Context, cmd interface{}) error {
	commandName := h.generateName(cmd)
	fmt.Printf("%s starting command processing: %s", time.Now().Format("02-01-06 15:04:05.999 Z0700"), commandName)

	err := h.handler.Handle(context, cmd)

	if err == nil {
		fmt.Printf("%s command processing - Success: %s", time.Now().Format("02-01-06 15:04:05.999 Z0700"), commandName)
	} else {
		fmt.Printf("%s command processing - Failure: %s", time.Now().Format("02-01-06 15:04:05.999 Z0700"), commandName)
	}

	return err
}
