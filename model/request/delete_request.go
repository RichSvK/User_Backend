package request

type DeleteUserRequest struct {
	UserId string `json:"user_id" validate:"required,uuid"`
}
