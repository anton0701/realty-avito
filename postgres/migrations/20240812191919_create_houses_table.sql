-- +goose Up
CREATE TABLE houses (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL,
    developer VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);

-- +goose Down
DROP TABLE IF EXISTS houses;
