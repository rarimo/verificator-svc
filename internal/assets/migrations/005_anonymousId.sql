-- +migrate Up
ALTER TABLE verify_users ADD COLUMN anonymous_id TEXT DEFAULT '';
ALTER TABLE verify_users ADD COLUMN nullifier TEXT DEFAULT '';

CREATE UNIQUE INDEX verify_users_anonymous_id_unique ON verify_users(anonymous_id) WHERE anonymous_id != '';
CREATE UNIQUE INDEX verify_users_nullifier_unique ON verify_users(nullifier) WHERE nullifier != '';
-- +migrate Down
ALTER TABLE verify_users DROP COLUMN anonymous_id;
ALTER TABLE verify_users DROP COLUMN nullifier;

DROP INDEX verify_users_anonymous_id_unique;
DROP INDEX verify_users_nullifier_unique;