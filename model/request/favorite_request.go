package request

import "stock_backend/model/entity"

type FavoriteUnderwriterRequest struct {
	Underwriter []entity.Underwriter `json:"underwriter" validate:"required" bson:"underwriter"`
}
