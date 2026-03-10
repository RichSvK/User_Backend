package request

type AddWatchlistRequest struct {
	Stock string `json:"stock" validate:"required"`
}
