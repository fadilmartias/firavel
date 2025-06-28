package requests

type UpdateProfileInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

type UpdatePasswordInput struct {
	CurrentPassword         string `json:"current_password" validate:"required,min=8,max=255"`
	NewPassword             string `json:"new_password" validate:"required,min=8,max=255"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,eqfield=NewPassword"`
}
