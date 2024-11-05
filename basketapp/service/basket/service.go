package basket

import (
	"context"
	"errors"
	"fmt"

	"git.gocasts.ir/remenu/beehive/types"

	"git.gocasts.ir/remenu/beehive/event"
)

var (
	ErrRestaurantMismatch  = errors.New("existing basket belongs to a different restaurant")
	ErrBasketNotRegistered = errors.New("the basket is already registered and cannot be modified")
)

// Repository defines the operations related to basket, supporting both Redis and PostgreSQL
type Repository interface {
	FindActiveBasket(ctx context.Context, userID types.ID) (Basket, error)
	AddItemToBasket(ctx context.Context, basketID types.ID, item BasketItem) error
	Create(ctx context.Context, basket Basket) (types.ID, error)
	Update(ctx context.Context, basket Basket) (types.ID, error)
	Delete(ctx context.Context, id types.ID) (bool, error)
	List(ctx context.Context) ([]Basket, error)
	CacheBasket(ctx context.Context, basket Basket) error
	GetCachedBasket(ctx context.Context, id types.ID) (Basket, error)
	
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
func (s Service) AddToBasket(ctx context.Context, basket Basket) (types.ID, error) {
	existingBasket, err := s.repository.FindActiveBasket(ctx, basket.UserID)
	if err != nil {
		return 0, fmt.Errorf("error retrieving basket: %v", err)
	}

	if existingBasket.ID != 0 {
		if existingBasket.RestaurantID != basket.RestaurantID {
			return 0, ErrRestaurantMismatch
		}

		if err := s.repository.AddItemToBasket(ctx, existingBasket.ID, basket.Items[0]); err != nil {
			return 0, fmt.Errorf("error adding item to existing basket: %v", err)
		}

		return existingBasket.ID, nil
	}

	newBasketID, err := s.repository.Create(ctx, basket)
	if err != nil {
		return 0, fmt.Errorf("error creating new basket: %v", err)
	}

	return newBasketID, nil

}

// UpdateBasket updates an existing basket
func (s Service) UpdateBasket(ctx context.Context, basket Basket) error {
	_, err := s.repository.Update(ctx, basket)
	if err != nil {
		return fmt.Errorf("error updating basket: %v", err)
	}
	return nil
}

// DeleteBasket deletes a basket by ID
func (s Service) DeleteBasket(ctx context.Context, id types.ID) error {
	_, err := s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting basket: %v", err)
	}
	return nil
}

// ListBaskets returns all baskets
func (s Service) ListBaskets(ctx context.Context) ([]Basket, error) {
	baskets, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing baskets: %v", err)
	}
	return baskets, nil
}

// GetBasketById retrieves a basket by its ID
func (s Service) GetBasketById(ctx context.Context, id types.ID) (Basket, error) {
	basket, err := s.repository.GetCachedBasket(ctx, id)
	if err == nil {
		// Basket found in cache
		return basket, nil
	}

	// If not found in cache, retrieve from PostgreSQL
	basketList, err := s.repository.List(ctx)
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
func (s Service) CacheBasket(ctx context.Context, basket Basket) error {
	err := s.repository.CacheBasket(ctx, basket)
	if err != nil {
		return fmt.Errorf("error caching basket: %v", err)
	}
	return nil
}

// GetCachedBasket retrieves a basket from Redis by its ID
func (s Service) GetCachedBasket(ctx context.Context, id types.ID) (Basket, error) {
	basket, err := s.repository.GetCachedBasket(ctx, id)
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
