-- +goose Up
ALTER TABLE flats ADD COLUMN moderator_id VARCHAR(255);

-- +goose Down
ALTER TABLE flats DROP COLUMN moderator_id;
