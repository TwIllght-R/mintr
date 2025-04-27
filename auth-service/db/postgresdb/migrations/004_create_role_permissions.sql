-- +goose Up

CREATE TABLE role_permissions (
  role_uuid UUID REFERENCES roles(uuid) ON DELETE CASCADE,
  permission_uuid UUID REFERENCES permissions(uuid) ON DELETE CASCADE,
  PRIMARY KEY (role_uuid, permission_uuid)
);

-- +goose Down
DROP TABLE role_permissions;