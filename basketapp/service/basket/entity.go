package basket

import (
	"time"

	"git.gocasts.ir/remenu/beehive/types"
)

// Basket represents the main basket entity
type Basket struct {
	ID             types.ID     `json:"id"`              // Basket unique ID
	UserID         types.ID     `json:"user_id"`         // ID of the user who owns the basket
	RestaurantID   types.ID     `json:"restaurant_id"`   // ID of the restaurant for the order
	ExpirationTime time.Time    `json:"expiration_time"` // Expiration date for the basket
	Items          []BasketItem `json:"items"`           // One-to-many relationship with BasketItem
	Status         BasketStatus `json:"status"`          // Status of the basket
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	IsDeleted      bool         `json:"is_deleted"`
}

// BasketItem represents the items in the basket
type BasketItem struct {
	ID          types.ID     `json:"id"`           // Item unique ID
	BasketID    types.ID     `json:"basket_id"`    // Reference to the parent basket
	FoodID      types.ID     `json:"food_id"`      // Food ID in the item
	FoodOptions []FoodOption `json:"food_options"` //List of food options
	FoodName    string       `json:"food_name"`    // Name of the product
	FoodPrice   types.Price  `json:"food_price"`   // Price of a single product
	Quantity    uint         `json:"quantity"`     // Quantity of the product in the basket
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// FoodOptionFoodOption represents a customizable option for a food item
type FoodOption struct {
	ID           types.ID    `json:"id"`             // Food Option unique ID
	BasketItemID types.ID    `json:"basket_item_id"` // Reference to the parent basket item
	OptionName   string      `json:"option_name"`    // Name of the food option (e.g., "Extra Cheese")
	OptionPrice  types.Price `json:"option_price"`   // Price for the option (e.g., 1000 Toman for extra cheese)
	Description  string      `json:"description"`    // Optional description for the option
}

// BasketStatus represents the possible statuses of a basket
type BasketStatus string

const (
	StatusNotRegistered         BasketStatus = "not_registered"
	StatusNotRegisteredCanceled BasketStatus = "not_registered_canceled"
	StatusRegistered            BasketStatus = "registered"
	StatusAccepted              BasketStatus = "accepted"
	StatusNotAccepted           BasketStatus = "not_accepted"
	StatusAcceptedCanceled      BasketStatus = "accepted_canceled"
	StatusAcceptedNotPaid       BasketStatus = "accepted_not_paid"
	StatusPaid                  BasketStatus = "paid"
)
