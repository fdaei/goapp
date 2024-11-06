package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"git.gocasts.ir/remenu/beehive/outbox"
	"git.gocasts.ir/remenu/beehive/types"

	"git.gocasts.ir/remenu/beehive/basketapp/service/basket"
	"github.com/redis/go-redis/v9"
)

// BasketRepo is the concrete implementation of the service.Repository interface
type BasketRepo struct {
	PostgreSQL *sql.DB       // PostgreSQL connection
	Redis      *redis.Client // Redis client connection
}

// NewBasketRepo creates a new instance of BasketRepo with PostgreSQL and Redis connections
func NewBasketRepo(db *sql.DB /*  redis *redis.Client */) basket.Repository {
	// TODO: add redis later
	return BasketRepo{
		PostgreSQL: db,
		//Redis:      redis,
	}
}

func (repo BasketRepo) FindActiveBasket(ctx context.Context, userID types.ID) (basket.Basket, error) {
	// TODO: what happen if we have more than one unregistered? Do I need to where clause for expiration time
	query := `
		SELECT id, user_id, restaurant_id, expiration_time, status, created_at, updated_at 
		FROM baskets 
		WHERE user_id = $1 AND status = $2 
		ORDER BY created_at DESC 
		LIMIT 1`

	row := repo.PostgreSQL.QueryRowContext(ctx, query, userID, basket.StatusNotRegistered)
	var b basket.Basket
	err := row.Scan(&b.ID, &b.UserID, &b.RestaurantID, &b.ExpirationTime, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return basket.Basket{}, nil
	} else if err != nil {
		return basket.Basket{}, fmt.Errorf("error retrieving active basket: %v", err)
	}

	return b, nil
}

func (repo BasketRepo) AddItemToBasket(ctx context.Context, basketID types.ID, item basket.BasketItem) error {
	query := `
		INSERT INTO basket_items (basket_id, food_id, food_name, food_price, quantity, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := repo.PostgreSQL.ExecContext(ctx, query, basketID, item.FoodID, item.FoodName, item.FoodPrice, item.Quantity, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("error adding item to basket: %v", err)
	}

	// Adding food options if available
	for _, option := range item.FoodOptions {
		optionQuery := `
			INSERT INTO food_options (basket_item_id, option_name, option_price, description) 
			VALUES ((SELECT id FROM basket_items WHERE basket_id=$1 AND food_id=$2 LIMIT 1), $3, $4, $5)`
		_, err := repo.PostgreSQL.ExecContext(ctx, optionQuery, basketID, item.FoodID, option.OptionName, option.OptionPrice, option.Description)
		if err != nil {
			return fmt.Errorf("error adding food option to basket item: %v", err)
		}
	}

	return nil
}

// Create inserts a new basket into PostgreSQL
func (repo BasketRepo) Create(ctx context.Context, b basket.Basket) (types.ID, error) {
	if b.Status == "" {
		b.Status = basket.StatusNotRegistered
	}

	// check if we have basket item to start transaction process
	if len(b.Items) > 0 {
		return repo.CreateWithTransaction(ctx, b)
	}

	// Start without transaction
	query := `
		INSERT INTO baskets (user_id, restaurant_id, expiration_time, created_at, updated_at, status, is_deleted) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`
	var basketID types.ID
	err := repo.PostgreSQL.QueryRowContext(
		ctx,
		query,
		b.UserID,
		b.RestaurantID,
		b.ExpirationTime,
		time.Now(),
		time.Now(),
		b.Status,
		b.IsDeleted,
	).Scan(&basketID)

	if err != nil {
		return 0, fmt.Errorf("error creating basket: %v", err)
	}

	return basketID, nil
}

// CreateWithTransaction creates a new basket with items and food options within a transaction.
func (repo BasketRepo) CreateWithTransaction(ctx context.Context, b basket.Basket) (types.ID, error) {
	// start transaction
	tx, err := repo.PostgreSQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `
	INSERT INTO baskets (user_id, restaurant_id, expiration_time, created_at, updated_at, status, is_deleted) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id`

	var basketID types.ID
	row := tx.QueryRowContext(
		ctx,
		query,
		b.UserID,
		b.RestaurantID,
		b.ExpirationTime,
		time.Now(),
		time.Now(),
		b.Status,
		b.IsDeleted,
	)

	// Check scan output
	if err := row.Scan(&basketID); err != nil {
		tx.Rollback()
		fmt.Printf("Error scanning basket ID: %v\n", err)
		return 0, fmt.Errorf("error scanning basket ID: %v", err)
	}

	fmt.Printf("Basket created with ID: %d\n", basketID)

	// add items to basket
	for _, item := range b.Items {
		if err := repo.addItemToBasketTx(ctx, tx, basketID, item); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("error adding item to basket: %v", err)
		}
	}

	// end of transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}

	return basketID, nil
}

// add items with transaction
func (repo BasketRepo) addItemToBasketTx(ctx context.Context, tx *sql.Tx, basketID types.ID, item basket.BasketItem) error {
	query := `
		INSERT INTO basket_items (basket_id, food_id, food_name, food_price, quantity, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`
	var itemID types.ID
	err := tx.QueryRowContext(ctx, query, basketID, item.FoodID, item.FoodName, item.FoodPrice, item.Quantity, time.Now(), time.Now()).Scan(&itemID)
	if err != nil {
		return err
	}

	// add options
	for _, option := range item.FoodOptions {
		if err := repo.addOptionToItemTx(ctx, tx, itemID, option); err != nil {
			return err
		}
	}

	return nil
}

// add options with transaction
func (repo BasketRepo) addOptionToItemTx(ctx context.Context, tx *sql.Tx, itemID types.ID, option basket.FoodOption) error {
	query := `
		INSERT INTO food_options (basket_item_id, option_name, option_price, description) 
		VALUES ($1, $2, $3, $4)`
	_, err := tx.ExecContext(ctx, query, itemID, option.OptionName, option.OptionPrice, option.Description)
	return err
}

// Update updates an existing basket in PostgreSQL
func (repo BasketRepo) Update(ctx context.Context, basket basket.Basket) (types.ID, error) {
	query := "UPDATE baskets SET restaurant_id=$1, expiration_time=$2, updated_at=$3 WHERE id=$4"
	_, err := repo.PostgreSQL.ExecContext(ctx, query, basket.RestaurantID, basket.ExpirationTime, basket.UpdatedAt, basket.ID)
	if err != nil {
		return 0, fmt.Errorf("error updating basket: %v", err)
	}
	return basket.ID, nil
}

// Delete removes a basket from PostgreSQL
func (repo BasketRepo) Delete(ctx context.Context, id types.ID) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM baskets WHERE id=$1)"
	err := repo.PostgreSQL.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking basket existence: %v", err)
	}

	if !exists {
		return false, nil
	}
	deleteQuery := "DELETE FROM baskets WHERE id=$1"
	_, err = repo.PostgreSQL.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return false, fmt.Errorf("error deleting basket: %v", err)
	}
	return true, nil
}

// List retrieves all baskets from PostgreSQL
func (repo BasketRepo) List(ctx context.Context) ([]basket.Basket, error) {
	query := "SELECT id, user_id, restaurant_id, expiration_time, created_at, updated_at FROM baskets"
	rows, err := repo.PostgreSQL.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving baskets: %v", err)
	}
	defer rows.Close()

	var baskets []basket.Basket
	for rows.Next() {
		var b basket.Basket
		if err := rows.Scan(&b.ID, &b.UserID, &b.RestaurantID, &b.ExpirationTime, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning basket: %v", err)
		}
		baskets = append(baskets, b)
	}

	return baskets, nil
}

// CacheBasket stores a basket in Redis
func (repo BasketRepo) CacheBasket(ctx context.Context, basket basket.Basket) error {
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
func (repo BasketRepo) GetCachedBasket(ctx context.Context, id types.ID) (basket.Basket, error) {
	var b basket.Basket
	basketData, err := repo.Redis.Get(context.Background(), fmt.Sprintf("basket:%d", id)).Result()
	if err != nil {
		return b, fmt.Errorf("error retrieving cached basket from Redis: %v", err)
	}
	if err := json.Unmarshal([]byte(basketData), &b); err != nil {
		return b, fmt.Errorf("error unmarshaling cached basket data: %v", err)
	}
	return b, nil
}

// TODO: should be more precise for specific use case like (updateing basket record and create outBoxEvent in a transaction).
// below function is just a sample just for creating OutBoxEvent and is incompelete
func (repo BasketRepo) CreateOutBox(ctx context.Context, outBoxEvent outbox.Event) (types.ID, error) {

	// TODO: use from sqlc
	var resultID types.ID
	query := "INSERT INTO outbox_events (id, topic, payload, is_published, retried_count, last_retried_at, published_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	// Use QueryRow to retrieve the resultID from RETURNING
	err := repo.PostgreSQL.QueryRowContext(
		ctx,
		query,
		outBoxEvent.ID,
		outBoxEvent.Topic,
		outBoxEvent.Payload,
		outBoxEvent.IsPublished,
		0,           // retried_count initialized to 0
		time.Time{}, // last_retried_at initialized as zero timestamp
		time.Time{}, // published_at initialized as zero timestamp
	).Scan(&resultID)

	if err != nil {
		return 0, fmt.Errorf("error creating outbox event: %v", err)
	}

	return resultID, nil
}
