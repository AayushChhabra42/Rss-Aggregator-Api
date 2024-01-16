-- +goose Up
Create Table feeds(
    id uuid Primary Key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name Text not null,
    url Text unique not null,
    user_id uuid not null references users(id) on delete cascade
);

-- +goose Down
Drop Table feeds;