package models

import "time"

type User struct {
	ID           string
	Name         string
	PasswordHash string
	Rights       int64
	CreatedAt    time.Time
}
