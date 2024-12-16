-- +migrate Up
ALTER TABLE verify_users ADD COLUMN expiration_lower_bound TEXT NOT NULL DEFAULT '0x303030303030';
-- +migrate Down
ALTER TABLE verify_users DROP COLUMN expiration_lower_bound;