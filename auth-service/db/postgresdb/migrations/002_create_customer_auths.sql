-- +goose Up

CREATE TABLE customer_auths (
  uuid UUID PRIMARY KEY NOT NULL,
  email VARCHAR UNIQUE NOT NULL,
  password TEXT NOT NULL,
  is_verified BOOLEAN DEFAULT FALSE,
  role_uuid UUID REFERENCES roles(uuid) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);


-- +goose Down
DROP TABLE customer_auths;

