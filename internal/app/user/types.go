package user

import (
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"
)

const (
	MALE Gender = iota
	FEMALE
	OTHER
)

type (
	Gender          int
	RegisterRequest struct {
		FirstName string `validate:"required" json:"first_name"`
		LastName  string `validate:"required" json:"last_name"`
		Gender    Gender `validate:"gte=0,lte=2" json:"gender"`
		Email     string `validate:"required,email" json:"email"`
		Password  string `validate:"required" json:"password"`
	}
	UpdateInfoRequest struct {
		FirstName string `validate:"required" json:"first_name"`
		LastName  string `validate:"required" json:"last_name"`
		Gender    Gender `validate:"gte=0,lte=2" json:"gender"`
	}

	ChangePasswordRequest struct {
		OldPassword string `validate:"required" json:"old_password"`
		NewPassword string `validate:"required" json:"new_password"`
	}

	User struct {
		ID        string    `bson:"_id"`
		FirstName string    `bson:"first_name"`
		LastName  string    `bson:"last_name"`
		Gender    Gender    `bson:"gender"`
		Email     string    `bson:"email"`
		Role      string    `bson:"role"`
		Password  string    `bson:"password"`
		Locked    bool      `bson:"locked"`
		CreatedAt time.Time `bson:"created_at"`
		UpdatedAt time.Time `bson:"updated_at"`
	}
)

func (u *User) ConvertUserToUserInfo() *types.UserInfo {

	return &types.UserInfo{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Gender:    int(u.Gender),
		Role:      u.Role,
		Locked:    u.Locked,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Password:  u.Password,
	}
}
