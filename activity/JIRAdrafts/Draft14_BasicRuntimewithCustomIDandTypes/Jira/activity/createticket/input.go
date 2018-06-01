package createticket

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func GetInputParameter(ivInput interface{}, dynamicMap map[string]interface{}) {
	activityLog.Info("Reading Input params")
	complexInput := ivInput.(*data.ComplexObject)

	metadata := make(map[string]interface{})
	json.Unmarshal([]byte(complexInput.Metadata), &metadata)

	props := metadata["properties"].(map[string]interface{})
	inputValues := complexInput.Value.(map[string]interface{})

	for key, value := range inputValues {
		customFields := props[key].(map[string]interface{})
		flogoJiraType := customFields["flogoJiraType"].(string)
		flogoJiraID := customFields["flogoJiraID"].(string)

		switch flogoJiraType {
		case "option":
			fieldValue := make(map[string]interface{})
			fieldValue["value"] = value
			dynamicMap[flogoJiraID] = fieldValue

		case "user":
			fieldValue := make(map[string]interface{})
			fieldValue["name"] = value
			dynamicMap[flogoJiraID] = fieldValue

		case "ArrayOfName":
			var arrayOfValue []interface{}
			for _, v1 := range value.([]interface{}) {
				fieldValue := make(map[string]interface{})
				fieldValue["name"] = v1
				arrayOfValue = append(arrayOfValue, fieldValue)
			}
			dynamicMap[flogoJiraID] = arrayOfValue

		case "string":
			dynamicMap[flogoJiraID] = value

		case "default":
			activityLog.Infof("JIRA datatype that is not handled by code :: %s ", flogoJiraType)
		}

	}

}
