-- +migrate Up
CREATE TYPE gender_enum AS ENUM ('', 'M', 'F', 'O');

-- +migrate StatementBegin
CREATE FUNCTION transform_gender(input TEXT) RETURNS gender_enum AS $$
BEGIN
    CASE UPPER(input)
        WHEN 'MALE' THEN RETURN 'M';
        WHEN 'FEMALE' THEN RETURN 'F';
        WHEN 'OTHERS' THEN RETURN 'O';
        WHEN '' THEN RETURN '';
        ELSE RAISE EXCEPTION 'invalid gender value: %', input;
    END CASE;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

ALTER TABLE verify_users ALTER COLUMN sex TYPE gender_enum USING transform_gender(sex);
DROP INDEX verify_users_nullifier_unique;
DROP INDEX verify_users_anonymous_id_unique;

-- +migrate Down
CREATE UNIQUE INDEX verify_users_anonymous_id_unique ON verify_users(anonymous_id) WHERE anonymous_id != '';
CREATE UNIQUE INDEX verify_users_nullifier_unique ON verify_users(nullifier) WHERE nullifier != '';
ALTER TABLE verify_users ALTER COLUMN sex TYPE TEXT;

DROP FUNCTION transform_gender;
DROP TYPE gender_enum;