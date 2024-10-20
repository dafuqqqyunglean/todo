package repository

import (
	_ "embed"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/jmoiron/sqlx"
)

type AuthorizationRepository interface {
	Create(user todo.User) (int, error)
	Get(username, password string) (todo.User, error)
}

type AuthorizationPostgres struct {
	db *sqlx.DB
}

func NewAuthorizationPostgres(db *sqlx.DB) *AuthorizationPostgres {
	return &AuthorizationPostgres{db: db}
}

//go:embed sql/CreateUser.sql
var createUser string

func (r *AuthorizationPostgres) Create(user todo.User) (int, error) {
	var id int

	row := r.db.QueryRow(createUser, user.Name, user.Username, user.Password) // stores information about the returned row from db
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

//go:embed sql/GetUser.sql
var getUser string

func (r *AuthorizationPostgres) Get(username, password string) (todo.User, error) {
	var user todo.User

	err := r.db.Get(&user, getUser, username, password)

	return user, err
}
