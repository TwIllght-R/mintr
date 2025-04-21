-- +goose Up

CREATE TABLE url_mappings (
  uuid UUID PRIMARY KEY NOT NULL,
  short_code VARCHAR(16) UNIQUE NOT NULL,
  original_url TEXT NOT NULL,
  title VARCHAR(50),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE url_mappings;