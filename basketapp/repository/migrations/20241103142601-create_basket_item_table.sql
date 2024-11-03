
-- +migrate Up
CREATE TABLE basket_items (
    id BIGINT PRIMARY KEY,
    basket_id BIGINT REFERENCES baskets(id) ON DELETE CASCADE,  
    food_id BIGINT NOT NULL,  
    food_name TEXT NOT NULL,  
    food_price BIGINT NOT NULL,  
    quantity INT DEFAULT 1, 
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- +migrate Down
DROP TABLE IF EXISTS basket_items;
