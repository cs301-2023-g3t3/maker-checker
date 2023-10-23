package utils

import (
	"fmt"
    // "strings"
)

func GetDifferences(oldData map[string]interface{}, updatedData map[string]interface{}) string {
    // for k, v := range oldData {
    //     if k != "id" {
    //         oldData[strings.ToLower(k)] = v
    //     }
    //     delete(oldData, k)
    // }
    fmt.Println(oldData)
    fmt.Println(updatedData)

    return "hello"
}
