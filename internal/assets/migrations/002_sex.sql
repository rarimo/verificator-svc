-- +migrate Up
ALTER TABLE verify_users ADD COLUMN sex BOOLEAN NOT NULL DEFAULT FALSE;
-- +migrate Down
ALTER TABLE verify_users DROP COLUMN sex;