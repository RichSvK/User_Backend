package entity

import (
	"github.com/google/uuid"
)

type Watchlist struct {
	UserId uuid.UUID `json:"user_id" validate:"required,uuid"`
	Stock  string    `json:"stock" validate:"required"`
}
