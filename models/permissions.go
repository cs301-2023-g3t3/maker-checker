package models

type Permission struct {
    Id              string                      `json:"_id" bson:"_id"`
    Route           string                      `json:"route" bson:"route" validation:"required"`
    Maker           []string                    `json:"maker" bson:"maker" validation:"required"`
    Checker         []string                    `json:"checker" bson:"checker" validation:"required"`
}
