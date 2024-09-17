-- +migrate Up
ALTER TABLE verify_users ADD COLUMN sex BOOLEAN;
-- +migrate Down
ALTER TABLE verify_users DROP COLUMN sex;