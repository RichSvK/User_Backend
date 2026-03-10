package response

type AddFavorite struct {
	Message string `json:"message"`
}

type RemoveFavorite struct {
	Message string `json:"message"`
}

type GetFavorites struct {
	Message string   `json:"message"`
	Data    []string `json:"data"`
}
