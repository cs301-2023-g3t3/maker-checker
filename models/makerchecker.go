package models

type MakercheckerList struct {
    Makercheckers []Makerchecker `json: "makercheckers" bson:",inline"`
}

type Makerchecker struct {
    Id          string      `json: "_id" bson: "_id"`
    MakerId     string      `json: "makerId" bson: "makerId"`
    CheckerId   string      `json: "checkerId" bson: "checkerId`
    Action      string      `json: "action" bson: "action"`
    Database    string      `json: "database" bson: "database"`
    Data        Data        `json: "data" bson: "data"`
}

type Data struct {
      User      User        `json:"user" bson:"user"`
      Point     Point       `json:"point" bson:"point"`
}

type User struct {
      Id        string      `json:"id" bson:"id"`
      Email     string      `json:"email" bson:"email"`
      FirstName string      `json:"firstName" bson:"firstName"`
      LastName  string      `json:"lastName" bson:"lastName"`
      Role      string      `json:"role" bson:"role"`
}

type Point struct {
      Id        string      `json:"id" bson:"id"`
      UserId    string      `json:"userId" bson:"userId"`
      Balance   float64     `json:"balance" bson:"balance"`
}
