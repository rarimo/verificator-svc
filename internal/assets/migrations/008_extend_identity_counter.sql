-- +migrate Up
ALTER TABLE verify_users 
    ALTER COLUMN identity_counter_upper_bound TYPE BIGINT,
    ALTER COLUMN identity_counter_lower_bound TYPE BIGINT,
    ALTER COLUMN identity_counter TYPE BIGINT;

-- +migrate Down  
ALTER TABLE verify_users 
    ALTER COLUMN identity_counter_upper_bound TYPE INTEGER,
    ALTER COLUMN identity_counter_lower_bound TYPE INTEGER,
    ALTER COLUMN identity_counter TYPE INTEGER;