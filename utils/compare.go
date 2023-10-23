package utils

import (
	"fmt"
)

func GetDifferences(oldData map[string]interface{}, updatedData map[string]interface{}) map[string]interface{} {
    var newValues = make(map[string]interface{})

    for k, v := range oldData {
        if k == "id" {
            newValues[k] = v
        }

        switch t := v.(type) {
        case string:
            if val, ok := updatedData[k].(string); ok {
                if (t != val) {
                    newValues[k] = val
                }
            }
            break
        case float64:
            if val, ok := updatedData[k].(float64); ok {
                if (t != val){
                    newValues[k] = val - t
                }
            }            
            break
        default:
            // TODO: handle data that isnt string or float64
            fmt.Println(t)
        }
    }

    return newValues
}
