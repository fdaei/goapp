-- +migrate Up
CREATE TYPE basket_status AS ENUM (
    'not_registered',
    'not_registered_canceled',
    'registered',
    'accepted',
    'not_accepted',
    'accepted_canceled',
    'accepted_not_paid',
    'paid'
);

-- +migrate Down
DROP TYPE IF EXISTS basket_status;
