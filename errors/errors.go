package errors

import "fmt"

type NotFoundError struct {
	id uint32
}

func NotFound(id uint32) *NotFoundError {
	return &NotFoundError{id: id}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Record with id: %d is not found!", e.id)
}

type UnexpectedError struct {
	Internal error
}

func Unexpected(err error) *UnexpectedError {
	return &UnexpectedError{Internal: err}
}

func (e *UnexpectedError) Error() string {
	return e.Internal.Error()
}
