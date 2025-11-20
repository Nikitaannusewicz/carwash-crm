-- +goose Up
ALTER TABLE locations RENAME COLUMN adress to address;
-- +goose Down
ALTER TABLE locations RENAME COLUMN adress to address;