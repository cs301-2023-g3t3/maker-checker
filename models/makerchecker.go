package models

type MakercheckerList struct {
    Makercheckers []Makerchecker `json:"makercheckers" bson:",inline"`
}

type Makerchecker struct {
    _Id             string      `json:"_id" bson:"_id"`
    MakercheckerId  string      `json:"makercheckerId" bson:"makercheckerId" validate:"required"`
    MakerId         string      `json:"makerId" bson:"makerId" validate:"required"`
    MakerEmail      string      `json:"makerEmail" bson:"makerEmail" validate:"required"`
    CheckerId       string      `json:"checkerId" bson:"checkerId" validate:"required"`
    CheckerEmail    string      `json:"checkerEmail" bson:"checkerEmail" validate:"required"`
    Action          string      `json:"action" bson:"action" validate:"required"`
    Status          string      `json:"status" bson:"status"`
    Database        string      `json:"database" bson:"database" validate:"required"`
    Data            *Data       `json:"data" bson:"data" validate:"required"`
}

type Data struct {
      User          *User       `json:"user,omitempty" bson:"user,omitempty"`
      Point         *Point      `json:"point,omitempty" bson:"point,omitempty"`
}

type User struct {
      Id            string      `json:"id" bson:"id" validate:"required"`
      Email         string      `json:"email,omitempty" bson:"email,omitempty"`
      FirstName     string      `json:"firstName,omitempty" bson:"firstName,omitempty"`
      LastName      string      `json:"lastName,omitempty" bson:"lastName,omitempty"`
      Role          string      `json:"role,omitempty" bson:"role,omitempty"`
}

type Point struct {
      Id            string      `json:"id" bson:"id" validate:"required"`
      UserId        string      `json:"userId,omitempty" bson:"userId,omitempty"`
      Balance       float64     `json:"balance,omitempty" bson:"balance,omitempty"`
}
