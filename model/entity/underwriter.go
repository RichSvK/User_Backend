package entity

type Underwriter struct {
	ID   string `json:"underwriter_id" bson:"underwriter_id"`
	Name string `json:"name" bson:"name"`
}
