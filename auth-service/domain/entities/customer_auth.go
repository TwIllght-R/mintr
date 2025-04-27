package entities

import "time"

type CustomerAuth struct {
	UUID       string
	Email      string
	Password   string
	IsVerified bool
	RoleUUID   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
