package models

type LoginModel struct {
	Email    *string `json:"email" validate:"email,required"`
	Password *string `json:"Password" validate:"required,min=6"`
}
