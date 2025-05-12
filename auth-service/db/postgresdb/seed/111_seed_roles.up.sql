-- +goose Up

INSERT INTO roles (uuid, name, created_at, updated_at) VALUES
  ('00000000-0000-0000-0000-000000000001', 'customer', now(), now()),
  ('00000000-0000-0000-0000-000000000002', 'moderator', now(), now()),
  ('00000000-0000-0000-0000-000000000003', 'admin', now(), now());

-- +goose Down

DELETE FROM roles WHERE uuid IN (
  '00000000-0000-0000-0000-000000000001',
  '00000000-0000-0000-0000-000000000002',
  '00000000-0000-0000-0000-000000000003'
);