package basketrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	basketmodel "git.gocasts.ir/remenu/beehive/service/basket/model"
	"github.com/redis/go-redis/v9"
)

// BasketRepository defines the operations related to basket, supporting both Redis and PostgreSQL
type Repository interface {
	Create(basket basketmodel.Basket) (sql.Result, error)
	Update(basket basketmodel.Basket) (sql.Result, error)
	Delete(id uint) (sql.Result, error)
	List() ([]basketmodel.Basket, error)
	CacheBasket(basket basketmodel.Basket) error
	GetCachedBasket(id uint) (basketmodel.Basket, error)
}

// BasketRepo is the concrete implementation of the BasketRepository
type BasketRepo struct {
	PostgreSQL *sql.DB       // PostgreSQL connection
	Redis      *redis.Client // Redis client connection
}

// NewBasketRepo creates a new instance of BasketRepo with PostgreSQL and Redis connections
func NewBasketRepo(db *sql.DB, redis *redis.Client) Repository {
	return &BasketRepo{
		PostgreSQL: db,
		Redis:      redis,
	}
}

// Create inserts a new basket into PostgreSQL
func (repo *BasketRepo) Create(basket basketmodel.Basket) (sql.Result, error) {

	// TODO: use from sqlc
	query := "INSERT INTO baskets (user_id, restaurant_id, expiration_time, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	result, err := repo.PostgreSQL.Exec(query, basket.UserID, basket.RestaurantID, basket.ExpirationTime, basket.CreatedAt, basket.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating basket: %v", err)
	}
	return result, nil
}

// Update updates an existing basket in PostgreSQL
func (repo *BasketRepo) Update(basket basketmodel.Basket) (sql.Result, error) {
	query := "UPDATE baskets SET restaurant_id=$1, expiration_time=$2, updated_at=$3 WHERE id=$4"
	result, err := repo.PostgreSQL.Exec(query, basket.RestaurantID, basket.ExpirationTime, basket.UpdatedAt, basket.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating basket: %v", err)
	}
	return result, nil
}

// Delete removes a basket from PostgreSQL
func (repo *BasketRepo) Delete(id uint) (sql.Result, error) {
	query := "DELETE FROM baskets WHERE id=$1"
	result, err := repo.PostgreSQL.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("error deleting basket: %v", err)
	}
	return result, nil
}

// List retrieves all baskets from PostgreSQL
func (repo *BasketRepo) List() ([]basketmodel.Basket, error) {
	query := "SELECT id, user_id, restaurant_id, expiration_time, created_at, updated_at FROM baskets"
	rows, err := repo.PostgreSQL.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving baskets: %v", err)
	}
	defer rows.Close()

	var baskets []basketmodel.Basket
	for rows.Next() {
		var b basketmodel.Basket
		if err := rows.Scan(&b.ID, &b.UserID, &b.RestaurantID, &b.ExpirationTime, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning basket: %v", err)
		}
		baskets = append(baskets, b)
	}

	return baskets, nil
}

// CacheBasket stores a basket in Redis
func (repo *BasketRepo) CacheBasket(basket basketmodel.Basket) error {
	basketData, err := json.Marshal(basket)
	if err != nil {
		return fmt.Errorf("error marshaling basket data for Redis: %v", err)
	}
	err = repo.Redis.Set(context.Background(), fmt.Sprintf("basket:%d", basket.ID), basketData, 0).Err()
	if err != nil {
		return fmt.Errorf("error caching basket in Redis: %v", err)
	}
	return nil
}

// GetCachedBasket retrieves a cached basket from Redis
func (repo *BasketRepo) GetCachedBasket(id uint) (basketmodel.Basket, error) {
	var b basketmodel.Basket
	basketData, err := repo.Redis.Get(context.Background(), fmt.Sprintf("basket:%d", id)).Result()
	if err != nil {
		return b, fmt.Errorf("error retrieving cached basket from Redis: %v", err)
	}
	if err := json.Unmarshal([]byte(basketData), &b); err != nil {
		return b, fmt.Errorf("error unmarshaling cached basket data: %v", err)
	}
	return b, nil
}
