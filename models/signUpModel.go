package models

type SignUp struct {
	First_name *string `json:"first_name" validate:"required,min=1,max=100"`
	Last_name  *string `json:"last_name" validate:"required,min=1,max=100"`
	Email      *string `json:"email"`
	Password   *string `json:"password"`
}
