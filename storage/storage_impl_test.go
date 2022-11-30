package storage

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"todo/errors"
)

func newTestStorage() (TodoStorage, func()) {
	filePrefix := uuid.New().String()
	dbFilePath := os.TempDir() + "/" + filePrefix + ".db"
	deleteDbFileIfExists(dbFilePath)

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		panic(err)
	}

	str := NewTodoStorage(db)
	err = str.Init()
	if err != nil {
		panic(err)
	}
	return str, func() {
		db.Close()
		deleteDbFileIfExists(dbFilePath)
	}
}

func deleteDbFileIfExists(file string) {
	if _, err := os.Stat(file); err == nil {
		// test db file exists
		err := os.Remove(file)
		if err != nil {
			panic(err)
		}
	}
}

func Test_todoStorage_Add(t1 *testing.T) {
	type args struct {
		todo AddTodo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Simple add", args: args{AddTodo{Name: "name1", Description: "description1"}}},
	}

	t, deleteFile := newTestStorage()
	defer deleteFile()
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			inserted := t.Add(tt.args.todo).MustGet().MustGet()
			assert.NotNil(t1, inserted.Id)
			assert.Equal(t1, tt.args.todo.Name, inserted.Name)
			assert.Equal(t1, tt.args.todo.Description, inserted.Description)
			assert.False(t1, inserted.Done)
		})
	}
}

func Test_todoStorage_MarkUndoneAndUndone(t1 *testing.T) {
	type args struct {
		create           *AddTodo
		updateDoneId     uint32
		updateDone       bool
		expected         *RecordTodo
		notFoundExpected bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Simple mark done",
			args: args{
				create: &AddTodo{
					Name:        "Name",
					Description: "Description",
				},
				updateDoneId: 0,
				updateDone:   true,
				expected: &RecordTodo{
					Id:          0,
					Name:        "Name",
					Description: "Description",
					Done:        true,
				},
			},
		}, {
			name: "Simple mark not done",
			args: args{
				create: &AddTodo{
					Name:        "Name",
					Description: "Description",
				},
				updateDoneId: 0,
				updateDone:   false,
				expected: &RecordTodo{
					Id:          0,
					Name:        "Name",
					Description: "Description",
					Done:        false,
				},
			},
		}, {
			name: "Simple mark not found",
			args: args{
				create:           nil,
				updateDoneId:     42,
				updateDone:       true,
				expected:         nil,
				notFoundExpected: true,
			},
		},
	}

	t, deleteFile := newTestStorage()
	defer deleteFile()
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if tt.args.create != nil {
				id := t.Add(*tt.args.create).MustGet().MustGet().Id
				tt.args.updateDoneId = id
				tt.args.expected.Id = id
			}

			var err error
			if tt.args.updateDone {
				err = t.MarkDone(tt.args.updateDoneId)
			} else {
				err = t.MarkUndone(tt.args.updateDoneId)
			}

			if tt.args.notFoundExpected {
				assert.Equal(t1, errors.NotFound(tt.args.updateDoneId), err)
			} else {
				assert.Equal(t1, t.FindById(tt.args.updateDoneId).MustGet().MustGet(), tt.args.expected)
			}
		})
	}
}

func Test_todoStorage_Update(t1 *testing.T) {
	type args struct {
		create   *AddTodo
		update   UpdateTodo
		expected *RecordTodo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Simple update", args: args{
				create: &AddTodo{
					Name:        "Old Name",
					Description: "Old Description",
				},
				update: UpdateTodo{
					Id:          0,
					Name:        "New Name",
					Description: "New Description",
				},
				expected: &RecordTodo{
					Id:          0,
					Name:        "New Name",
					Description: "New Description",
					Done:        false,
				}},
		}, {
			name: "Update not found", args: args{
				create: nil,
				update: UpdateTodo{
					Id:          123,
					Name:        "New Name",
					Description: "New Description",
				},
				expected: nil},
		},
	}

	t, deleteFile := newTestStorage()
	defer deleteFile()
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if tt.args.create != nil {
				res := t.Add(*tt.args.create)
				id := res.MustGet().MustGet().Id
				tt.args.update.Id = id

				if tt.args.expected != nil {
					tt.args.expected.Id = id
				}
			}

			updated := t.Update(tt.args.update)
			if tt.args.expected != nil {
				assert.Equal(t1, tt.args.expected, updated.MustGet().MustGet())
			} else {
				assert.Equal(t1, errors.NotFound(tt.args.update.Id), updated.Error())
			}
		})
	}
}
