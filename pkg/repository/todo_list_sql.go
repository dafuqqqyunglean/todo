package repository

import (
	_ "embed"
	"fmt"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoListRepository interface {
	Create(ctx utility.AppContext, userId int, list todo.TodoList) (int, error)
	GetAll(ctx utility.AppContext, userId int) ([]todo.TodoList, error)
	GetById(ctx utility.AppContext, userId, listId int) (todo.TodoList, error)
	Delete(ctx utility.AppContext, userId, listId int) error
	Update(ctx utility.AppContext, userId, listId int, input todo.UpdateListInput) error
}

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

//go:embed sql/CreateList.sql
var createList string

//go:embed sql/CreateUsersLists.sql
var createUsersLists string

func (r *TodoListPostgres) Create(ctx utility.AppContext, userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.BeginTx(ctx.Ctx, nil)
	if err != nil {
		return 0, err
	}

	var id int
	row := tx.QueryRow(createList, list.Title, list.Description) // stores information about the returned row from db
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec(createUsersLists, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

//go:embed sql/GetAllLists.sql
var getAllLists string

func (r *TodoListPostgres) GetAll(ctx utility.AppContext, userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList

	err := r.db.SelectContext(ctx.Ctx, &lists, getAllLists, userId)

	return lists, err
}

//go:embed sql/GetListById.sql
var getListById string

func (r *TodoListPostgres) GetById(ctx utility.AppContext, userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	err := r.db.GetContext(ctx.Ctx, &list, getListById, userId, listId)

	return list, err
}

//go:embed sql/DeleteList.sql
var deleteList string

func (r *TodoListPostgres) Delete(ctx utility.AppContext, userId, listId int) error {
	_, err := r.db.ExecContext(ctx.Ctx, deleteList, userId, listId)

	return err
}

//go:embed sql/UpdateList.sql
var updateList string

func (r *TodoListPostgres) Update(ctx utility.AppContext, userId, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(updateList, setQuery, argId, argId+1)
	args = append(args, listId, userId)

	_, err := r.db.ExecContext(ctx.Ctx, query, args...)
	return err
}
