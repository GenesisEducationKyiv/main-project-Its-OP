package application

type IValidator[T any] interface {
	Validate(T) error
}
