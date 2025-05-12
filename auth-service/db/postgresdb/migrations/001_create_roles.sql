-- +goose Up

CREATE TABLE roles (
  uuid UUID PRIMARY KEY NOT NULL,
  name VARCHAR UNIQUE NOT NULL,   -- เช่น "admin", "user", "moderator"
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down

DROP TABLE roles;