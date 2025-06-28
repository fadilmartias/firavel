package requests

type RegisterInput struct {
	Name                 string `json:"name" validate:"required,max=255"`
	Email                string `json:"email" validate:"required,email,max=255"`
	Phone                string `json:"phone" validate:"required,max=255"`
	Password             string `json:"password" validate:"required,min=8,max=255"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type LoginInput struct {
	Credential string `json:"credential" validate:"required,max=255"`
	Password   string `json:"password" validate:"required,min=8,max=255"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type ResetPasswordInput struct {
	Token                string `json:"token" validate:"required,max=255"`
	Password             string `json:"password" validate:"required,min=8,max=255"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type VerifyEmailInput struct {
	Token string `json:"token" validate:"required,max=255"`
}
