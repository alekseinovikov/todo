package service

import . "github.com/samber/mo"

type CreateTodo struct {
	Name    string
	Content string
}

type UpdateTodo struct {
	Id          uint32
	Name        string
	Description string
}

type Todo struct {
	Id          uint32
	Name        string
	Description string
	Done        bool
}

type TodoService interface {
	Save(ct CreateTodo) (*Todo, error)
	Update(ut UpdateTodo) (*Todo, error)
	FindById(id uint32) (Option[*Todo], error)
	MarkDone(id uint32) error
	MarkUndone(id uint32) error
}
