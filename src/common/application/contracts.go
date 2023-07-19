//go:generate mockery --name ILogger
package application

import "context"

type IValidator[T any] interface {
	Validate(T) error
}

type ILogger interface {
	LogInformation(message string) error
	LogDebug(message string) error
	LogError(err error, message string) error
}

type ICommandBus interface {
	Send(ctx context.Context, cmd any) error
}
