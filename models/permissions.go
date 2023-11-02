package models

type Permission struct {
    Id              string                      `json:"_id" bson:"_id"`
    Endpoint        string                      `json:"endpoint" bson:"endpoint" validation:"required"`
    Maker           []float64                   `json:"maker" bson:"maker" validation:"required"`
    Checker         []float64                    `json:"checker" bson:"checker" validation:"required"`
}
