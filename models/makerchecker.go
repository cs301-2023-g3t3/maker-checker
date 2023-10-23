package models

type MakercheckerList struct {
    Makercheckers []Makerchecker `json:"makercheckers" bson:",inline"`
}

type Makerchecker struct {
    _Id             string                      `json:"_id" bson:"_id"`
    MakercheckerId  string                      `json:"makercheckerId" bson:"makercheckerId"`
    MakerId         string                      `json:"makerId" bson:"makerId" validate:"required"`
    MakerEmail      string                      `json:"makerEmail" bson:"makerEmail" validate:"required"`
    CheckerId       string                      `json:"checkerId" bson:"checkerId" validate:"required"`
    CheckerEmail    string                      `json:"checkerEmail" bson:"checkerEmail" validate:"required"`
    Action          string                      `json:"action" bson:"action" validate:"required"`
    Database        string                      `json:"database" bson:"database" validate:"required"`
    Status          string                      `json:"status" bson:"status"`
    Data            map[string]interface{}      `json:"data" bson:"data" validate:"required"`
}
