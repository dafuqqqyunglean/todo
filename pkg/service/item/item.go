package item

import (
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository/cache"
)

type TodoItemService interface {
	Create(ctx utility.AppContext, userId, listId int, item todo.TodoItem) (int, error)
	GetAll(ctx utility.AppContext, userId, listId int) ([]todo.TodoItem, error)
	GetById(ctx utility.AppContext, userId, itemId int) (todo.TodoItem, error)
	Delete(ctx utility.AppContext, userId, itemId int) error
	Update(ctx utility.AppContext, userId, itemId int, input todo.UpdateItemInput) error
}

type ImplTodoItem struct {
	repo     repository.TodoItemRepository
	listRepo repository.TodoListRepository
	cache    cache.RedisCache
}

func NewTodoItemService(repo repository.TodoItemRepository, listRepo repository.TodoListRepository, cache cache.RedisCache) *ImplTodoItem {
	return &ImplTodoItem{
		repo:     repo,
		listRepo: listRepo,
		cache:    cache,
	}
}

func (s *ImplTodoItem) Create(ctx utility.AppContext, userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(ctx, userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, listId, item)
}

func (s *ImplTodoItem) GetAll(ctx utility.AppContext, userId, listId int) ([]todo.TodoItem, error) {
	items, err := s.repo.GetAll(ctx, userId, listId)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ImplTodoItem) GetById(ctx utility.AppContext, userId, itemId int) (todo.TodoItem, error) {
	item, err := s.cache.GetItem(ctx, userId, itemId)
	if err == nil {
		return item, nil
	}

	item, err = s.repo.GetById(ctx, userId, itemId)
	if err != nil {
		return item, err
	}

	s.cache.SetItem(ctx, userId, itemId, item)

	return item, nil
}

func (s *ImplTodoItem) Delete(ctx utility.AppContext, userId, itemId int) error {
	err := s.repo.Delete(ctx, userId, itemId)
	if err != nil {
		return err
	}
	s.cache.Delete(ctx, userId, itemId)

	return nil
}

func (s *ImplTodoItem) Update(ctx utility.AppContext, userId, itemId int, input todo.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	err := s.repo.Update(ctx, userId, itemId, input)
	if err != nil {
		return err
	}

	s.cache.Delete(ctx, userId, itemId)
	return nil
}
