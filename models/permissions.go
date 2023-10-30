package models

type Permission struct {
    Id              string                      `json:"_id" bson:"_id"`
    Endpoint        string                      `json:"endpoint" bson:"endpoint" validation:"required"`
    Maker           []string                    `json:"maker" bson:"maker" validation:"required"`
    Checker         []string                    `json:"checker" bson:"checker" validation:"required"`
}
