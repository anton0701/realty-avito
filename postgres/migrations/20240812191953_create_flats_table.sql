-- +goose Up
CREATE TABLE flats (
    id SERIAL PRIMARY KEY,
    house_id INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    price INTEGER NOT NULL,
    rooms INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);

-- +goose Down
DROP TABLE IF EXISTS flats;