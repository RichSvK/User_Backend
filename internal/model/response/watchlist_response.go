package response

type RemoveWatchlistResponse struct {
	Message string `json:"message"`
}

type AddWatchlistResponse struct {
	Message string `json:"message"`
}

type GetWatchlistResponse struct {
	Message string   `json:"message"`
	Stocks  []string `json:"stocks"`
}
