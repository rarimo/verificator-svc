-- +migrate Up
-- +migrate Up
ALTER TABLE verify_users RENAME COLUMN sex TO sex_enable;
ALTER TABLE verify_users ADD COLUMN sex TEXT;
-- +migrate Down
ALTER TABLE verify_users DROP COLUMN sex;
ALTER TABLE verify_users DROP COLUMN sex_enable;