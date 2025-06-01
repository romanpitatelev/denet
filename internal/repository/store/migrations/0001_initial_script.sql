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
    deleted_at TIMESTAMP WITH TIME ZONE,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_user_id ON users(user_id);

CREATE TABLE tasks
(
    id uuid PRIMARY KEY,
    user_id REFERENCES users (id) NOT NULL,
    type VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
}

-- +migrate Down
DROP INDEX IF EXISTS idx_user_id;
