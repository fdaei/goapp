package basket

import (
	"fmt"

	"git.gocasts.ir/remenu/beehive/event"
	basketmodel "git.gocasts.ir/remenu/beehive/service/basket/model"
	basketrepo "git.gocasts.ir/remenu/beehive/service/basket/repository"
)

// Service defines the operations related to basket
type Service interface {
	CreateBasket(basket basketmodel.Basket) (uint, error)
	UpdateBasket(basket basketmodel.Basket) error
	DeleteBasket(id uint) error
	ListBaskets() ([]basketmodel.Basket, error)
	GetBasketById(id uint) (basketmodel.Basket, error)
	GetBasketItemsById(id uint) ([]basketmodel.BasketItem, error)
	CacheBasket(basket basketmodel.Basket) error
	GetCachedBasket(id uint) (basketmodel.Basket, error)
	PurchaseSucceedHandler(event event.Event) error
	PurchaseFailedHandler(event event.Event) error 
}

// BasketService is the concrete implementation of Service
type BasketService struct {
	BasketRepo basketrepo.Repository
}

// NewBasketService creates a new instance of BasketService
func NewBasketService(repo basketrepo.Repository) Service {
	return &BasketService{
		BasketRepo: repo,
	}
}

// CreateBasket creates a new basket
func (s *BasketService) CreateBasket(basket basketmodel.Basket) (uint, error) {
	result, err := s.BasketRepo.Create(basket)
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
func (s *BasketService) UpdateBasket(basket basketmodel.Basket) error {
	_, err := s.BasketRepo.Update(basket)
	if err != nil {
		return fmt.Errorf("error updating basket: %v", err)
	}
	return nil
}

// DeleteBasket deletes a basket by ID
func (s *BasketService) DeleteBasket(id uint) error {
	_, err := s.BasketRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting basket: %v", err)
	}
	return nil
}

// ListBaskets returns all baskets
func (s *BasketService) ListBaskets() ([]basketmodel.Basket, error) {
	baskets, err := s.BasketRepo.List()
	if err != nil {
		return nil, fmt.Errorf("error listing baskets: %v", err)
	}
	return baskets, nil
}

// GetBasketById retrieves a basket by its ID
func (s *BasketService) GetBasketById(id uint) (basketmodel.Basket, error) {
	basket, err := s.BasketRepo.GetCachedBasket(id)
	if err == nil {
		// Basket found in cache
		return basket, nil
	}

	// If not found in cache, retrieve from PostgreSQL
	basketList, err := s.BasketRepo.List()
	if err != nil {
		return basketmodel.Basket{}, fmt.Errorf("error retrieving basket: %v", err)
	}
	for _, b := range basketList {
		if b.ID == id {
			return b, nil
		}
	}

	return basketmodel.Basket{}, fmt.Errorf("basket not found")
}

func (s *BasketService) GetBasketItemsById(id uint) ([]basketmodel.BasketItem, error) {
	// TODO: complete this section
	return []basketmodel.BasketItem{}, nil
}

// CacheBasket caches a basket in Redis
func (s *BasketService) CacheBasket(basket basketmodel.Basket) error {
	err := s.BasketRepo.CacheBasket(basket)
	if err != nil {
		return fmt.Errorf("error caching basket: %v", err)
	}
	return nil
}

// GetCachedBasket retrieves a basket from Redis by its ID
func (s *BasketService) GetCachedBasket(id uint) (basketmodel.Basket, error) {
	basket, err := s.BasketRepo.GetCachedBasket(id)
	if err != nil {
		return basketmodel.Basket{}, fmt.Errorf("error retrieving cached basket: %v", err)
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
