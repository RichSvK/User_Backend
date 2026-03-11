package response

type AddFavoriteResponse struct {
	Message string `json:"message"`
}

type RemoveFavoriteResponse struct {
	Message string `json:"message"`
}

type GetFavoritesResponse struct {
	Message string   `json:"message"`
	Data    []string `json:"data"`
}
