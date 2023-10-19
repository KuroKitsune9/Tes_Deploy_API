package model

import (
	"time"
)

type LoginRequest struct {
	Email    string `form:"email" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type UserModel struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Umur      int64      `json:"umur"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Password  string     `json:"-"`
}
