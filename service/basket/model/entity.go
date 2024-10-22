package basketmodel

import "time"

// Basket represents the main basket entity
type Basket struct {
	ID             uint         `json:"id"`              // Basket unique ID
	UserID         uint         `json:"user_id"`         // ID of the user who owns the basket
	RestaurantID   uint         `json:"restaurant_id"`   // ID of the restaurant for the order
	ExpirationTime time.Time    `json:"expiration_time"` // Expiration date for the basket
	Items          []BasketItem `json:"items"`           // One-to-many relationship with BasketItem
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// BasketItem represents the items in the basket
type BasketItem struct {
	ID          uint         `json:"id"`           // Item unique ID
	BasketID    uint         `json:"basket_id"`    // Reference to the parent basket
	FoodID      uint         `json:"food_id"`      // Food ID in the item
	FoodOptions []FoodOption `json:"food_options"` //List of food options
	FoodName    string       `json:"food_name"`    // Name of the product
	FoodPrice   uint         `json:"food_price"`   // Price of a single product
	Quantity    uint         `json:"quantity"`     // Quantity of the product in the basket
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// FoodOption represents a customizable option for a food item
type FoodOption struct {
	OptionName  string `json:"option_name"`  // Name of the food option (e.g., "Extra Cheese")
	OptionPrice uint   `json:"option_price"` // Price for the option (e.g., 1000 Toman for extra cheese)
	Description string `json:"description"`  // Optional description for the option
}
