-- +migrate Up
CREATE TABLE users
(
    id uuid PRIMARY KEY,
    name VARCHAR NOT NULL,
    email VARCHAR,
    role VARCHAR,
    points NUMERIC NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX users_unique_email_and_deleted_at_null_idx 
    ON users (email)
    WHERE deleted_at IS NULL;

CREATE TABLE tasks
(
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES users (id) NOT NULL,
    type VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reference
(
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES users (id) NOT NULL,
    reference_id uuid REFERENCES users (id) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP INDEX IF EXISTS users_unique_email_and_deleted_at_null_idx;
DROP TABLE IF EXISTS reference;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;
