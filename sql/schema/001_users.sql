-- +goose Up
Create Table users(
    id uuid Primary Key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name Text not null
);

-- +goose Down
Drop Table users;