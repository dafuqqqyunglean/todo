package list

import (
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository/cache"
)

type TodoListService interface {
	Create(ctx utility.AppContext, userId int, list todo.TodoList) (int, error)
	GetAll(ctx utility.AppContext, userId int) ([]todo.TodoList, error)
	GetById(ctx utility.AppContext, userId, listId int) (todo.TodoList, error)
	Delete(ctx utility.AppContext, userId, listId int) error
	Update(ctx utility.AppContext, userId, listId int, input todo.UpdateListInput) error
}

type ImplTodoList struct {
	repo  repository.TodoListRepository
	cache cache.RedisCache
}

func NewTodoListService(repo repository.TodoListRepository, cache cache.RedisCache) *ImplTodoList {
	return &ImplTodoList{
		repo:  repo,
		cache: cache,
	}
}

func (s *ImplTodoList) Create(ctx utility.AppContext, userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(ctx, userId, list)
}

func (s *ImplTodoList) GetAll(ctx utility.AppContext, userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(ctx, userId)
}

func (s *ImplTodoList) GetById(ctx utility.AppContext, userId, listId int) (todo.TodoList, error) {
	list, err := s.cache.GetList(ctx, userId, listId)
	if err == nil {
		return list, nil
	}

	list, err = s.repo.GetById(ctx, userId, listId)
	if err != nil {
		return list, err
	}

	s.cache.SetList(ctx, userId, listId, list)

	return list, nil
}

func (s *ImplTodoList) Delete(ctx utility.AppContext, userId, listId int) error {
	err := s.repo.Delete(ctx, userId, listId)
	if err != nil {
		return err
	}
	s.cache.Delete(ctx, userId, listId)

	return nil
}

func (s *ImplTodoList) Update(ctx utility.AppContext, userId, listId int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	err := s.repo.Update(ctx, userId, listId, input)
	if err != nil {
		return err
	}

	s.cache.Delete(ctx, userId, listId)
	return nil
}
