package types

import (
	"time"
)

const (
	MALE   Gender = 1
	FEMALE Gender = 2
	OTHER  Gender = 3
)

const (
	RoleAdmin  Role = "admin"
	RoleNormal Role = "normal"
)

type (
	Gender   int
	UserInfo struct {
		ID        string    `json:"_id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Gender    int       `json:"gender"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		Password  string    `json:"password,omitempty"`
		Locked    bool      `json:"locked"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	Role string
)
