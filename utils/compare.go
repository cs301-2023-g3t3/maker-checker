package utils

import (
    "fmt"
)

func GetDifferences(action string, oldData map[string]interface{}, updatedData map[string]interface{}) string {
    fmt.Println(action)
    fmt.Println(oldData)
    fmt.Println(updatedData)

    return "hello"
}
