-- +goose Up
CREATE TABLE permissions (
  uuid UUID PRIMARY KEY NOT NULL,
  name VARCHAR UNIQUE NOT NULL,  -- เช่น "url:create", "url:view", "analytics:view"
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE permissions;