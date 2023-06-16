package domain

type DataConsistencyError struct {
	Message string
}

func (e *DataConsistencyError) Error() string {
	return e.Message
}

type EndpointInaccessibleError struct {
	Message string
}

func (e *EndpointInaccessibleError) Error() string {
	return e.Message
}

const databaseErrorMessage = "There was an issue with the database. Please try again later."

type DatabaseError struct {
	NestedError error
}

func (e *DatabaseError) Error() string {
	return databaseErrorMessage
}

func (e *DatabaseError) Unwrap() error {
	return e.NestedError
}
