-- +migrate Up
ALTER TABLE verify_users ADD COLUMN nationality_enable BOOLEAN NOT NULL DEFAULT FALSE;
-- +migrate Down
ALTER TABLE verify_users DROP COLUMN nationality_enable;