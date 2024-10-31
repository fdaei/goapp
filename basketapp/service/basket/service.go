package basket

import (
	"database/sql"
	"fmt"
	"git.gocasts.ir/remenu/beehive/types"

	"git.gocasts.ir/remenu/beehive/event"
)

// Repository defines the operations related to basket, supporting both Redis and PostgreSQL
type Repository interface {
	Create(basket Basket) (uint64, error)
	Update(basket Basket) (sql.Result, error)
	Delete(id uint) (sql.Result, error)
	List() ([]Basket, error)
	CacheBasket(basket Basket) error
	GetCachedBasket(id types.ID) (Basket, error)
}

// Service is the concrete implementation of Service
type Service struct {
	repository Repository
}

// NewService creates a new instance of Service
func NewService(repo Repository) Service {
	return Service{
		repository: repo,
	}
}

// CreateBasket creates a new basket
func (s Service) CreateBasket(basket Basket) (uint, error) {
	id, err := s.repository.Create(basket)
	if err != nil {
		return 0, fmt.Errorf("error creating basket: %v", err)
	}

	if err != nil {
		return 0, fmt.Errorf("error retrieving last insert ID: %v", err)
	}
	return uint(id), nil
}

// UpdateBasket updates an existing basket
func (s Service) UpdateBasket(basket Basket) error {
	_, err := s.repository.Update(basket)
	if err != nil {
		return fmt.Errorf("error updating basket: %v", err)
	}
	return nil
}

// DeleteBasket deletes a basket by ID
func (s Service) DeleteBasket(id uint) error {
	_, err := s.repository.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting basket: %v", err)
	}
	return nil
}

// ListBaskets returns all baskets
func (s Service) ListBaskets() ([]Basket, error) {
	baskets, err := s.repository.List()
	if err != nil {
		return nil, fmt.Errorf("error listing baskets: %v", err)
	}
	return baskets, nil
}

// GetBasketById retrieves a basket by its ID
func (s Service) GetBasketById(id types.ID) (Basket, error) {
	basket, err := s.repository.GetCachedBasket(id)
	if err == nil {
		// Basket found in cache
		return basket, nil
	}

	// If not found in cache, retrieve from PostgreSQL
	basketList, err := s.repository.List()
	if err != nil {
		return Basket{}, fmt.Errorf("error retrieving basket: %v", err)
	}
	for _, b := range basketList {
		if b.ID == id {
			return b, nil
		}
	}

	return Basket{}, fmt.Errorf("basket not found")
}

func (s Service) GetBasketItemsById(id uint) ([]BasketItem, error) {
	// TODO: complete this section
	return []BasketItem{}, nil
}

// CacheBasket caches a basket in Redis
func (s Service) CacheBasket(basket Basket) error {
	err := s.repository.CacheBasket(basket)
	if err != nil {
		return fmt.Errorf("error caching basket: %v", err)
	}
	return nil
}

// GetCachedBasket retrieves a basket from Redis by its ID
func (s Service) GetCachedBasket(id types.ID) (Basket, error) {
	basket, err := s.repository.GetCachedBasket(id)
	if err != nil {
		return Basket{}, fmt.Errorf("error retrieving cached basket: %v", err)
	}
	return basket, nil
}

func (s Service) PurchaseSucceedHandler(event event.Event) error {
	fmt.Println("PurchaseSucceedHandler")
	// TODO: add transaction to create and outbox Event Message and publish notify event
	return nil
}

func (s Service) PurchaseFailedHandler(event event.Event) error {
	fmt.Println("PurchaseFailedHandler")
	return nil
}
