package utils

import (
	"fmt"
)

func GetDifferences(oldData map[string]interface{}, updatedData map[string]interface{}) (int, map[string]interface{}) {
    var newValues = make(map[string]interface{})

    for k, v := range oldData {
        if k == "id" {
            newValues[k] = v
        }

        if _, found := updatedData[k]; !found {
            msg := fmt.Sprintf("Key '%s' is not found in new data", k)
            return 400, map[string]interface{}{"data": msg}
        }

        switch t := v.(type) {
        case string:
            if val, ok := updatedData[k].(string); ok {
                newValues[k] = val
            } else {
                msg := fmt.Sprintf("Key '%s' with value '%v' is not string type, but %T", k, updatedData[k], updatedData[k])
                return 400, map[string]interface{}{"data": msg}
            }
            break
        case float64:
            if val, ok := updatedData[k].(float64); ok {
                if (t != val){
                    newValues[k] = val - t
                } else {
                    newValues[k] = t
                }
            } else {
                msg := fmt.Sprintf("Key '%s' with value '%v' is not float64 type, but %T", k, updatedData[k], updatedData[k])
                return 400, map[string]interface{}{"data": msg}
            }
            break
        default:
            msg := fmt.Sprintf("Key '%s' with value '%v' is not string or float64 type, but type %T", k, updatedData[k], updatedData[k])
            return 400, map[string]interface{}{"data": msg}
        }
    }

    return 200, newValues
}
