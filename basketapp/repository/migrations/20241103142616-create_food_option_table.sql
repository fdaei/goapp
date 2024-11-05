
-- +migrate Up
CREATE TABLE food_options (
    id SERIAL PRIMARY KEY,
    basket_item_id BIGINT REFERENCES basket_items(id) ON DELETE CASCADE,  
    option_name TEXT NOT NULL,
    option_price NUMERIC(10, 2) DEFAULT 0, 
    description TEXT  
);
-- +migrate Down
DROP TABLE IF EXISTS food_options;
