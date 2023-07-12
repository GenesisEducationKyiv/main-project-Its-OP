package application

type IValidator[T any] interface {
	Validate(T) error
}

type ILogger interface {
	LogInformation(message string) error
	LogDebug(message string) error
	LogError(err error, message string) error
}
