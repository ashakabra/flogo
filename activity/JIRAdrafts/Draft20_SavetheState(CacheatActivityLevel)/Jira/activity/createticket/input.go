package createticket

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

const (
	FIELDS    = "fields"
	PROJECT   = "project"
	ISSUETYPE = "issuetype"

	KEY   = "key"
	NAME  = "name"
	VALUE = "value"
)

func CreateInputFields(outputMap map[string]interface{}, projectKey string, issueType string, input interface{}) map[string]interface{} {
	inputJSON := make(map[string]interface{})
	dynamicMap := make(map[string]interface{})

	dynamicMap[PROJECT] = createKeyValue(KEY, projectKey)
	dynamicMap[ISSUETYPE] = createKeyValue(NAME, issueType)

	GetInputParameter(input, dynamicMap, outputMap)

	inputJSON[FIELDS] = dynamicMap
	return inputJSON
}

func createKeyValue(key, value string) map[string]string {
	m := make(map[string]string)
	m[key] = value
	return m
}

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
			fieldValue[VALUE] = value
			dynamicMap[flogoJiraID] = fieldValue

		case "user", "priority", "group", "version":
			fieldValue := make(map[string]interface{})
			fieldValue[NAME] = value
			dynamicMap[flogoJiraID] = fieldValue

		case "project":
			fieldValue := make(map[string]interface{})
			fieldValue[KEY] = value
			dynamicMap[flogoJiraID] = fieldValue

		case "ArrayOfName":
			var arrayOfName []interface{}
			for _, v1 := range value.([]interface{}) {
				fieldValue := make(map[string]interface{})
				fieldValue[NAME] = v1
				arrayOfName = append(arrayOfName, fieldValue)
			}
			dynamicMap[flogoJiraID] = arrayOfName

		case "ArrayOfValue":
			var arrayOfValue []interface{}
			for _, v1 := range value.([]interface{}) {
				fieldValue := make(map[string]interface{})
				fieldValue[VALUE] = v1
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
