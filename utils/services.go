package utils

import "makerchecker/models"

func ProcessMicroserviceTypes(makerchecker models.Makerchecker) (string, string) {
    switch str := makerchecker.Database; str {
    case "users":
        return "user-storage-api", "users"
    case "points":
        return "points-ledger-api", "points"
    default:
        return "Error", "Error"
    }
}
