//go:generate mockery --name ILogger
package application

import "context"

type IValidator[T any] interface {
	Validate(T) error
}

type ILogger interface {
	Info(message string) error
	Debug(message string) error
	Error(err error, message string) error
}

type ICommandBus interface {
	Send(ctx context.Context, cmd any) error
}
