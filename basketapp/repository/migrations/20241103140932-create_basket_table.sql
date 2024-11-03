
-- +migrate Up
CREATE TABLE baskets (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,  
    restaurant_id BIGINT NOT NULL,  
    expiration_time TIMESTAMPTZ DEFAULT (NOW() + INTERVAL '30 minutes'),  
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE
);

-- +migrate Down
DROP TABLE IF EXISTS baskets;