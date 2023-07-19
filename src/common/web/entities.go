package web

type Response[T any] struct {
	Code         int
	Body         *T
	ErrorMessage string
	Successful   bool
}
