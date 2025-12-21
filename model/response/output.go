package response

import "time"

type Output struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
	Data    any       `json:"data"`
}

type FailedResponse struct {
	Message string `json:"message"`
}
