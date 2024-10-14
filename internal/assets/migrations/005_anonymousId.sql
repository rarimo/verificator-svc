-- +migrate Up
ALTER TABLE verify_users ADD COLUMN anonymous_id TEXT UNIQUE;
ALTER TABLE verify_users ADD COLUMN nullifier TEXT UNIQUE;
-- +migrate Down
ALTER TABLE verify_users DROP COLUMN anonymous_id;
ALTER TABLE verify_users DROP COLUMN nullifier;