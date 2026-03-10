package response

type RegisterResponse struct {
	Message string `json:"message"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type VerifyResponse struct {
	Message string `json:"message"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type UserProfileResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
