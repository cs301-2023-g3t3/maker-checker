package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type MakercheckerList struct {
    Makercheckers []Makerchecker `json: "makercheckers" bson:",inline"`
}

type Makerchecker struct {
    _Id              primitive.ObjectID    `json:"_id"`
    MakerId         string      `json:"makerId" validate:"required"`
    MakerEmail      string      `json:"makerEmail" validate:"required"`
    CheckerId       string      `json:"checkerId" validate:"required"`
    CheckerEmail    string      `json:"checkerEmail" validate:"required"`
    Action          string      `json:"action" validate:"required"`
    Database        string      `json:"database" validate:"required"`
    Data            *Data        `json:"data" validate:"required"`
}

type Data struct {
      User          *User        `json:"user,omitempty"`
      Point         *Point       `json:"point,omitempty"`
}

type User struct {
      Id            string      `json:"id" validate:"required"`
      Email         string      `json:"email"`
      FirstName     string      `json:"firstName"`
      LastName      string      `json:"lastName"`
      Role          string      `json:"role"`
}

type Point struct {
      Id            string      `json:"id" validate:"required"`
      UserId        string      `json:"userId"`
      Balance       float64     `json:"balance"`
}
