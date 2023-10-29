package utils

func ProcessMicroserviceTypes(endpoint string) (string, string) {
    switch endpoint {
    case "users":
        return "user-storage-api", "users/accounts"
    case "points":
        return "points-ledger-api", "points/accounts"
    default:
        return "Error", "Error"
    }
}
