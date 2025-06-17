-- +migrate Up
ALTER TABLE verify_users 
    ADD COLUMN birth_date_lower_bound TEXT,
    ADD COLUMN birth_date_upper_bound TEXT,
    ADD COLUMN event_data TEXT,
    ADD COLUMN expiration_date_upper_bound TEXT,
    ADD COLUMN identity_counter INTEGER DEFAULT 0,
    ADD COLUMN identity_counter_lower_bound INTEGER DEFAULT 0,
    ADD COLUMN identity_counter_upper_bound INTEGER DEFAULT 0,
    ADD COLUMN selector TEXT NOT NULL DEFAULT '',
    ADD COLUMN timestamp_lower_bound TEXT,
    ADD COLUMN timestamp_upper_bound TEXT;

-- +migrate Down
ALTER TABLE verify_users 
    DROP COLUMN birth_date_lower_bound,
    DROP COLUMN birth_date_upper_bound,
    DROP COLUMN event_data,
    DROP COLUMN expiration_date_upper_bound,
    DROP COLUMN identity_counter,
    DROP COLUMN identity_counter_lower_bound,
    DROP COLUMN identity_counter_upper_bound,
    DROP COLUMN selector,
    DROP COLUMN timestamp_lower_bound,
    DROP COLUMN timestamp_upper_bound;