package domain

type DataConsistencyError struct {
	Message string
}

func (e DataConsistencyError) Error() string {
	return e.Message
}

type EndpointInaccessibleError struct {
	Message string
}

func (e EndpointInaccessibleError) Error() string {
	return e.Message
}

const internalErrorMessage = "Internal server error. Please try again later."

type InternalError struct {
	NestedError error
}

func (e InternalError) Error() string {
	return internalErrorMessage
}

func (e InternalError) Unwrap() error {
	return e.NestedError
}
