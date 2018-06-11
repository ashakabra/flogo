package createticket

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func GetInputParameter(ivInput interface{}, dynamicMap map[string]interface{}, outputMap map[string]interface{}) {
	activityLog.Info("Reading Input params")
	complexInput := ivInput.(*data.ComplexObject)

	metadata := make(map[string]interface{})
	json.Unmarshal([]byte(complexInput.Metadata), &metadata)

	props := metadata["properties"].(map[string]interface{})
	inputValues := complexInput.Value.(map[string]interface{})

	for key, value := range inputValues {
		outputMap[key] = value
		customFields := props[key].(map[string]interface{})
		flogoJiraType := customFields["flogoJiraType"].(string)
		flogoJiraID := customFields["flogoJiraID"].(string)

		switch flogoJiraType {
		case "option":
			fieldValue := make(map[string]interface{})
			fieldValue["value"] = value
			dynamicMap[flogoJiraID] = fieldValue

		case "user", "priority", "group", "version":
			fieldValue := make(map[string]interface{})
			fieldValue["name"] = value
			dynamicMap[flogoJiraID] = fieldValue

		case "project":
			fieldValue := make(map[string]interface{})
			fieldValue["key"] = value
			dynamicMap[flogoJiraID] = fieldValue

		case "ArrayOfName":
			var arrayOfName []interface{}
			for _, v1 := range value.([]interface{}) {
				fieldValue := make(map[string]interface{})
				fieldValue["name"] = v1
				arrayOfName = append(arrayOfName, fieldValue)
			}
			dynamicMap[flogoJiraID] = arrayOfName

		case "ArrayOfValue":
			var arrayOfValue []interface{}
			for _, v1 := range value.([]interface{}) {
				fieldValue := make(map[string]interface{})
				fieldValue["value"] = v1
				arrayOfValue = append(arrayOfValue, fieldValue)
			}
			dynamicMap[flogoJiraID] = arrayOfValue

		case "ArrayOfString":
			var arrayOfString []string
			for _, v1 := range value.([]interface{}) {
				arrayOfString = append(arrayOfString, v1.(string))
			}
			dynamicMap[flogoJiraID] = arrayOfString

		case "string", "number", "date", "datetime", "any":
			dynamicMap[flogoJiraID] = value

		case "default":
			activityLog.Infof("JIRA datatype that is not handled by code :: %s ", flogoJiraType)
		}

	}

}
