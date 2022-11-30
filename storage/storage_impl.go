package storage

import (
	"database/sql"
	"fmt"
	. "github.com/samber/mo"
	"log"
	"todo/errors"
)

type todoStorage struct {
	db *sql.DB
}

func NewTodoStorage(db *sql.DB) TodoStorage {
	return &todoStorage{
		db: db,
	}
}

func (t *todoStorage) Init() error {
	sqlStmt := `
		create table if not exists todos(id integer not null primary key autoincrement, 
										name text not null, 
										description text null default null, 
										done integer default 0);
	`

	_, err := t.db.Exec(sqlStmt)
	return err
}

func (t *todoStorage) Close() error {
	return t.db.Close()
}

func (t *todoStorage) Add(todo AddTodo) (*RecordTodo, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, errors.Unexpected(err)
	}

	stmt, err := tx.Prepare(`insert into todos(name, description) 
									VALUES (?,?)`)
	if err != nil {
		return nil, errors.Unexpected(err)
	}
	defer closeInternal(stmt)

	result, err := stmt.Exec(todo.Name, todo.Description)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.Unexpected(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Unexpected(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Unexpected(err)
	}

	found, err := t.FindById(uint32(id))
	if err != nil {
		return nil, err
	}

	if found.IsAbsent() {
		return nil, errors.Unexpected(fmt.Errorf("can't find entity after insert"))
	}

	return found.MustGet(), nil
}

func (t *todoStorage) FindById(id uint32) (Option[*RecordTodo], error) {
	tx, err := t.db.Begin()
	if err != nil {
		return None[*RecordTodo](), errors.Unexpected(err)
	}

	rows, err := tx.Query("SELECT id, name, description, done FROM todos WHERE id = ? LIMIT 1", id)
	if err != nil {
		return None[*RecordTodo](), errors.Unexpected(err)
	}
	defer closeInternal(rows)

	if !rows.Next() {
		return noneResult, nil
	}

	var rec RecordTodo
	var done int
	err = rows.Scan(&rec.Id, &rec.Name, &rec.Description, &done)
	if err != nil {
		_ = tx.Rollback()
		return None[*RecordTodo](), errors.Unexpected(err)
	}

	err = tx.Commit()
	if err != nil {
		return None[*RecordTodo](), errors.Unexpected(err)
	}

	rec.Done = done > 0
	return Some(&rec), nil
}

func (t *todoStorage) Update(todo UpdateTodo) (*RecordTodo, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, errors.Unexpected(err)
	}

	result, err := tx.Exec("UPDATE todos SET name = ?, description = ? WHERE id = ?", todo.Name, todo.Description, todo.Id)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.Unexpected(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Unexpected(err)
	}

	if count, err := result.RowsAffected(); count < 1 || err != nil {
		return nil, errors.NotFound(todo.Id)
	}

	found, err := t.FindById(todo.Id)
	if err != nil {
		return nil, err
	}

	if found.IsAbsent() {
		return nil, fmt.Errorf("can't find entity after update")
	}

	return found.MustGet(), nil
}

func (t *todoStorage) MarkDone(id uint32) error {
	return t.updateDone(id, 1)
}

func (t *todoStorage) MarkUndone(id uint32) error {
	return t.updateDone(id, 0)
}

func (t *todoStorage) updateDone(id uint32, done int) error {
	tx, err := t.db.Begin()
	if err != nil {
		return errors.Unexpected(err)
	}

	result, err := tx.Exec("UPDATE todos SET done = ? WHERE id = ?", done, id)
	if err != nil {
		_ = tx.Rollback()
		return errors.Unexpected(err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Unexpected(err)
	}

	if count, _ := result.RowsAffected(); count < 1 || err != nil {
		return errors.NotFound(id)
	}

	return nil
}

type Closer interface {
	Close() error
}

func closeInternal(closer Closer) {
	err := closer.Close()
	if err != nil {
		log.Fatal(err)
	}
}

var noneResult = None[*RecordTodo]()
