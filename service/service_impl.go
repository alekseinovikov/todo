package service

import (
	"todo/storage"

	. "github.com/samber/mo"
)

func NewTodoService(storage storage.TodoStorage) TodoService {
	return &todoService{storage: storage}
}

type todoService struct {
	storage storage.TodoStorage
}

func (t todoService) Save(ct CreateTodo) (*Todo, error) {
	result, err := t.storage.Add(ct.toStorage())
	if err != nil {
		return nil, err
	}
	return fromRecord(result), nil
}

func (t todoService) Update(ut UpdateTodo) (*Todo, error) {
	result, err := t.storage.Update(ut.toStorage())
	if err != nil {
		return nil, err
	}

	return fromRecord(result), nil
}

func (t todoService) FindById(id uint32) (Option[*Todo], error) {
	result, err := t.storage.FindById(id)
	if err != nil {
		return None[*Todo](), err
	}

	return fromOptionRecord(result), nil
}

func (t todoService) MarkDone(id uint32) error {
	return t.storage.MarkDone(id)
}

func (t todoService) MarkUndone(id uint32) error {
	return t.storage.MarkUndone(id)
}

func (c *CreateTodo) toStorage() storage.AddTodo {
	return storage.AddTodo{
		Name:        c.Name,
		Description: c.Content,
	}
}

func (u *UpdateTodo) toStorage() storage.UpdateTodo {
	return storage.UpdateTodo{
		Id:          u.Id,
		Name:        u.Name,
		Description: u.Description,
	}
}

func fromRecord(r *storage.RecordTodo) *Todo {
	return &Todo{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Done:        r.Done,
	}
}

func fromOptionRecord(opt Option[*storage.RecordTodo]) Option[*Todo] {
	if opt.IsAbsent() {
		return None[*Todo]()
	}

	return Some(fromRecord(opt.MustGet()))
}
