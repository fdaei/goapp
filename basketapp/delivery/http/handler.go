package http

import (
	"errors"
	"net/http"
	"reflect"

	"git.gocasts.ir/remenu/beehive/basketapp/service/basket"
	"git.gocasts.ir/remenu/beehive/types"
	"github.com/go-playground/validator/v10"
	echo "github.com/labstack/echo/v4"
)

// AddToBasketRequest represents the expected request body for adding a food item to basket
type AddToBasketRequest struct {
	RestaurantID types.ID            `json:"restaurant_id" validate:"required"`
	UserID       types.ID            `json:"user_id" validate:"required"`
	FoodID       types.ID            `json:"food_id" validate:"required"`
	FoodName     string              `json:"food_name" validate:"required"`
	Quantity     uint                `json:"quantity" validate:"required,min=1"`
	FoodPrice    types.Price         `json:"food_price" validate:"required"`
	FoodOptions  []FoodOptionRequest `json:"food_options"` // Optional
}

// FoodOption represents the options for a food item
type FoodOptionRequest struct {
	OptionName  string      `json:"option_name"`
	OptionPrice types.Price `json:"option_price"`
	Description string      `json:"description"`
}

type Handler struct {
	BasketService basket.Service
	Validator     *validator.Validate
}

// NewHandler initializes a new Handler with the basket service
func NewHandler(basketService basket.Service) Handler {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag := fld.Tag.Get("json")
		if jsonTag == "-" {
			return ""
		}
		if jsonTag == "" {
			return fld.Name
		}
		return jsonTag
	})

	return Handler{
		BasketService: basketService,
		Validator:     validate,
	}
}

// AddToBasket handles adding a food item to the basket
func (h Handler) AddToBasket(c echo.Context) error {
	var req AddToBasketRequest

	// Parsing the request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	if err := h.Validator.Struct(req); err != nil {
		var invalidFields []string
		for _, err := range err.(validator.ValidationErrors) {
			invalidFields = append(invalidFields, err.Field())
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":                   "validation failed: need to send all required field",
			"missing_required_fields": invalidFields,
		})
	}

	basketItem := basket.BasketItem{
		FoodID:      req.FoodID,
		FoodName:    req.FoodName,
		FoodPrice:   req.FoodPrice,
		Quantity:    req.Quantity,
		FoodOptions: convertFoodOptions(req.FoodOptions),
	}

	currentBasket := basket.Basket{
		UserID:       req.UserID,
		RestaurantID: req.RestaurantID,
		Items:        []basket.BasketItem{basketItem},
	}

	ctx := c.Request().Context()
	basketID, err := h.BasketService.AddToBasket(ctx, currentBasket)
	// Handling specific errors from the service layer
	if err != nil {
		if errors.Is(err, basket.ErrRestaurantMismatch) {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "please clear the current basket before adding items from a different restaurant",
			})
		}
		if errors.Is(err, basket.ErrBasketNotRegistered) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "the basket is already registered and cannot be modified",
			})
		}
		// General error handling for unexpected errors
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not add item to basket",
		})
	}

	// Return success response with Basket ID
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Item added to basket successfully",
		"basket_id": basketID,
	})
}

func convertFoodOptions(options []FoodOptionRequest) []basket.FoodOption {
	var foodOptions []basket.FoodOption
	for _, option := range options {
		foodOptions = append(foodOptions, basket.FoodOption{
			OptionName:  option.OptionName,
			OptionPrice: option.OptionPrice,
			Description: option.Description,
		})
	}
	return foodOptions
}
