package request

type AddFavoriteUnderwriterRequest struct {
	UnderwriterId string `json:"underwriter_id" validate:"required"`
}
