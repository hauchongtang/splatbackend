package models

type TaskResult struct {
	First_name  *string `json:"first_name" validate:"required,min=1,max=100"`
	Last_name   *string `json:"last_name" validate:"required,min=1,max=100"`
	Task_name   *string `json:"taskName"`
	Module_code *string `json:"moduleCode"`
	Duration    *string `json:"duration"`
	Hidden      bool    `json:"hidden"`
	User_id     string  `json:"user_id"`
}
