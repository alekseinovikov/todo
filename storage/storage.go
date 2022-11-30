package storage

import . "github.com/samber/mo"

type AddTodo struct {
	Name        string
	Description string
}

type UpdateTodo struct {
	Id          uint32
	Name        string
	Description string
}

type RecordTodo struct {
	Id          uint32
	Name        string
	Description string
	Done        bool
}

type TodoStorage interface {
	Init() error
	Close() error
	Add(todo AddTodo) (*RecordTodo, error)
	FindById(id uint32) (Option[*RecordTodo], error)
	Update(todo UpdateTodo) (*RecordTodo, error)
	MarkDone(id uint32) error
	MarkUndone(id uint32) error
}
