-- +migrate Up

CREATE TABLE verify_users
(
    user_id         TEXT PRIMARY KEY NOT NULL,
    user_id_hash    TEXT             NOT NULL,
    age_lower_bound INT,
    nationality     TEXT,
    created_at      TIMESTAMP        NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    uniqueness      BOOLEAN          NOT NULL,
    event_id        TEXT,
    status          TEXT             NOT NULL,
    proof           JSON             NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS verify_users;
