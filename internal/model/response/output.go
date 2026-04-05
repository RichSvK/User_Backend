package response

type WebResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type FailedResponse struct {
	Message string `json:"message"`
}
