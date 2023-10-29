package models

type Makerchecker struct {
    Id             string                       `json:"_id" bson:"_id"`
    MakercheckerId  string                      `json:"makercheckerId" bson:"makercheckerId"`
    MakerId         string                      `json:"makerId" bson:"makerId"`
    MakerEmail      string                      `json:"makerEmail" bson:"makerEmail"`
    CheckerId       string                      `json:"checkerId" bson:"checkerId"`
    CheckerEmail    string                      `json:"checkerEmail" bson:"checkerEmail"`
    Action          string                      `json:"action" bson:"action"`
    Database        string                      `json:"database" bson:"database"`
    Status          string                      `json:"status" bson:"status"`
    Data            map[string]interface{}      `json:"data" bson:"data"`
}
