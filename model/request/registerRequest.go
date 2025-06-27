package request

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
	Username string `json:"username" validate:"required,min=3,max=30"`
}
