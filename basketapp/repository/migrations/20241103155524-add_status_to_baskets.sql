-- +migrate Up
ALTER TABLE baskets
ADD COLUMN status basket_status DEFAULT 'not_registered';

-- +migrate Down
ALTER TABLE baskets
DROP COLUMN IF EXISTS status;
