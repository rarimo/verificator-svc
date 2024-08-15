-- +migrate Up

CREATE TABLE verify_users
(
    user_id       TEXT PRIMARY KEY NOT NULL,
    user_id_hash  TEXT NOT NULL,
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    status        BOOLEAN     NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS verify_users;
