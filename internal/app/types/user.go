package types

import (
	"time"
)

const (
	MALE Gender = iota
	FEMALE
	OTHER
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
		Password  string    `json:"password"`
		Locked    bool      `json:"locked"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
