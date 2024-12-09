-- +migrate Up
ALTER TABLE verify_users ADD COLUMN expiration_lower_bound TEXT NOT NULL DEFAULT '52983525027888';
-- +migrate Down
ALTER TABLE verify_users DROP COLUMN expiration_lower_bound;