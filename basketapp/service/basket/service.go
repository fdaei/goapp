package basket

import (
	"database/sql"
	"fmt"

	"git.gocasts.ir/remenu/beehive/event"
)

// BasketRepository defines the operations related to basket, supporting both Redis and PostgreSQL
type Repository interface {
	Create(basket Basket) (sql.Result, error)
	Update(basket Basket) (sql.Result, error)
	Delete(id uint) (sql.Result, error)
	List() ([]Basket, error)
	CacheBasket(basket Basket) error
	GetCachedBasket(id uint) (Basket, error)
}

// Service defines the operations related to basket
type Service interface {
	CreateBasket(basket Basket) (uint, error)
	UpdateBasket(basket Basket) error
	DeleteBasket(id uint) error
	ListBaskets() ([]Basket, error)
	GetBasketById(id uint) (Basket, error)
	GetBasketItemsById(id uint) ([]BasketItem, error)
	CacheBasket(basket Basket) error
	GetCachedBasket(id uint) (Basket, error)
	PurchaseSucceedHandler(event event.Event) error
	PurchaseFailedHandler(event event.Event) error
}

// BasketService is the concrete implementation of Service
type BasketService struct {
	repository Repository
}

// NewBasketService creates a new instance of BasketService
func NewBasketService(repo Repository) Service {
	return &BasketService{
		repository: repo,
	}
}

// CreateBasket creates a new basket
func (s *BasketService) CreateBasket(basket Basket) (uint, error) {
	result, err := s.repository.Create(basket)
	if err != nil {
		return 0, fmt.Errorf("error creating basket: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error retrieving last insert ID: %v", err)
	}
	return uint(id), nil
}

// UpdateBasket updates an existing basket
func (s *BasketService) UpdateBasket(basket Basket) error {
	_, err := s.repository.Update(basket)
	if err != nil {
		return fmt.Errorf("error updating basket: %v", err)
	}
	return nil
}

// DeleteBasket deletes a basket by ID
func (s *BasketService) DeleteBasket(id uint) error {
	_, err := s.repository.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting basket: %v", err)
	}
	return nil
}

// ListBaskets returns all baskets
func (s *BasketService) ListBaskets() ([]Basket, error) {
	baskets, err := s.repository.List()
	if err != nil {
		return nil, fmt.Errorf("error listing baskets: %v", err)
	}
	return baskets, nil
}

// GetBasketById retrieves a basket by its ID
func (s *BasketService) GetBasketById(id uint) (Basket, error) {
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

func (s *BasketService) GetBasketItemsById(id uint) ([]BasketItem, error) {
	// TODO: complete this section
	return []BasketItem{}, nil
}

// CacheBasket caches a basket in Redis
func (s *BasketService) CacheBasket(basket Basket) error {
	err := s.repository.CacheBasket(basket)
	if err != nil {
		return fmt.Errorf("error caching basket: %v", err)
	}
	return nil
}

// GetCachedBasket retrieves a basket from Redis by its ID
func (s *BasketService) GetCachedBasket(id uint) (Basket, error) {
	basket, err := s.repository.GetCachedBasket(id)
	if err != nil {
		return Basket{}, fmt.Errorf("error retrieving cached basket: %v", err)
	}
	return basket, nil
}

func (s *BasketService) PurchaseSucceedHandler(event event.Event) error {
	fmt.Println("PurchaseSucceedHandler")
	return nil
}

func (s *BasketService) PurchaseFailedHandler(event event.Event) error {
	fmt.Println("PurchaseFailedHandler")
	return nil
}
