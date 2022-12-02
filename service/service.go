package service

import . "github.com/samber/mo"

type CreateTodo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateTodo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Todo struct {
	Id          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type TodoService interface {
	Save(ct CreateTodo) (*Todo, error)
	Update(id uint32, ut UpdateTodo) (*Todo, error)
	FindById(id uint32) (Option[*Todo], error)
	MarkDone(id uint32) error
	MarkUndone(id uint32) error
}
