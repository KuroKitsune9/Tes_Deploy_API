package model

import (
	"time"
)

type RegisReponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Umur      string    `json:"umur"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisRequest struct {
	Name     string `form:"name" validate:"required"`
	Email    string `form:"email" validate:"required,email"`
	Umur     int    `form:"umur" validate:"required,numeric"`
	Password string `form:"password" validate:"required"`
}
