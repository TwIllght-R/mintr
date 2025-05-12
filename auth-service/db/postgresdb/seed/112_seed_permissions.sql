-- +goose Up

INSERT INTO permissions (uuid, name, created_at, updated_at) VALUES
  ('10000000-0000-0000-0000-000000000001', 'shorten_url', now(), now()),
  ('10000000-0000-0000-0000-000000000002', 'view_url_stats', now(), now()),
  ('10000000-0000-0000-0000-000000000003', 'create_custom_alias', now(), now()),
  ('10000000-0000-0000-0000-000000000004', 'ban_user', now(), now());

-- +goose Down
DELETE FROM permissions WHERE uuid IN (
  '10000000-0000-0000-0000-000000000001',
  '10000000-0000-0000-0000-000000000002',
  '10000000-0000-0000-0000-000000000003',
  '10000000-0000-0000-0000-000000000004'
);