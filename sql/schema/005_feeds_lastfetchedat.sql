-- +goose Up
Alter table feeds add column last_fetched_at timestamp;
-- +goose Down
Alter table feeds drop column last_fetched_at;