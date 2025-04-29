-- +goose Up

CREATE TABLE url_mappings (
  uuid UUID PRIMARY KEY NOT NULL,
  owner_uuid UUID NOT NULL,
  short_code VARCHAR(16) UNIQUE NOT NULL,
  original_url TEXT NOT NULL,
  title VARCHAR(50),
  utm_source VARCHAR(50),
  utm_medium VARCHAR(50),
  utm_campaign VARCHAR(50),
  utm_term VARCHAR(50),
  utm_content VARCHAR(50),
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE url_mappings;