package entities

type Permission struct {
	UUID string
	Name string
}

type RolePermission struct {
	RoleUUID       string
	PermissionUUID string
}

type Role struct {
	UUID string
	Name string
}
