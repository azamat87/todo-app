package repository

import (
	"errors"
	todoapp "golang_ninja/todo-app"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTodoItemPostgrs_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item   todoapp.TodoItem
	}

	type mockBehavier func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavier mockBehavier
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				listId: 1,
				item: todoapp.TodoItem{
					Title:       "test title",
					Description: "test Description",
				},
			},
			id: 2,
			mockBehavier: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_item").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectQuery("INSERT INTO list_items").WithArgs(args.listId, id).WillReturnRows(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

			},
		},
		{
			name: "Empty Fields",
			args: args{
				listId: 1,
				item: todoapp.TodoItem{
					Title:       "",
					Description: "test Description",
				},
			},
			mockBehavier: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("fields error"))
				mock.ExpectQuery("INSERT INTO todo_item").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectRollback()

			},
			wantErr: true,
		},
		{
			name: "Second Insert error",
			args: args{
				listId: 1,
				item: todoapp.TodoItem{
					Title:       "test title",
					Description: "test Description",
				},
			},
			id: 2,
			mockBehavier: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_item").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectQuery("INSERT INTO list_items").WithArgs(args.listId, id).WillReturnError(errors.New("Second Insert error"))
				mock.ExpectRollback()

			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavier(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.listId, testCase.args.item)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}

}
