package errors

type NotFoundError struct {}

func (e *NotFoundError) Error() string {
	return "resource not found"
}

func NotFound() error {
	return &NotFoundError{}
}
