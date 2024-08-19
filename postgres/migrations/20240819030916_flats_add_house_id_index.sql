-- +goose Up
CREATE INDEX idx_house_id ON flats(house_id);

-- +goose Down
DROP INDEX IF EXISTS idx_house_id;