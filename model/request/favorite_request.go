package request

type FavoriteUnderwriterRequest struct {
	UnderwriterId string `json:"underwriter_id" validate:"required"`
}
