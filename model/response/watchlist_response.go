package response

type WatchlistResponse struct {
	Message string `json:"message"`
}

type GetWatchlistResponse struct {
	Message string   `json:"message"`
	Stocks  []string `json:"stocks"`
}
