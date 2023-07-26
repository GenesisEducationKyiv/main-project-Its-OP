//go:generate mockery --name ILogger
package application

import "context"

type IValidator[T any] interface {
	Validate(*T) error
}

type ILogger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type ICommandBus interface {
	Send(ctx context.Context, cmd any) error
}
