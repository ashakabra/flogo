package createticket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var activityLog = logger.GetLogger("jira-activity-createticket")

const (
	ivConnection = "Connection"
	ivProject    = "project"
	ivIssueType  = "issueType"
	ivInput      = "input"
	ovOutput     = "output"
)

// CreateTicketActivity struct
type CreateTicketActivity struct {
	metadata *activity.Metadata
}

// NewActivity to create a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &CreateTicketActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *CreateTicketActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval runtime execution
func (a *CreateTicketActivity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Infof("Executing Jira Create Ticket")

	connector := context.GetInput(ivConnection)
	if connector == nil || len(connector.(map[string]interface{})) == 0 {
		return false, activity.NewError("Jira connection is not configured", "JIRA-CREATETICKET-4001", nil)
	}

	//Read connection details
	conn := GetConnection(connector)

	//create json input
	outputMap := make(map[string]interface{})
	inputJSON := CreateInputFields(outputMap, context.GetInput(ivProject).(string), context.GetInput(ivIssueType).(string), context.GetInput(ivInput))

	//call create rest api
	response, err := CreateTicket(conn, inputJSON)
	if err != nil {
		return false, activity.NewError(fmt.Sprintf("Failed with error %s ", err.Error()), "JIRA-CREATETICKET-4002", nil)
	} else {
		jsonResponseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false, activity.NewError(fmt.Sprintf("Error reading json, %s", err.Error()), "JIRA-CREATETICKET-4003", nil)
		}
		defer response.Body.Close()
		if response.StatusCode >= 400 {
			return false, activity.NewError(fmt.Sprintf("Failed with error %s", jsonResponseData), "JIRA-CREATETICKET-4004", nil)
		}

		queryResponse := make(map[string]interface{})
		err = json.Unmarshal(jsonResponseData, &queryResponse)

		if err != nil {
			return false, activity.NewError(fmt.Sprintf("Unmarshal json err %s", err.Error()), "JIRA-CREATETICKET-4005", nil)
		}

		outputMap["IssueID"] = queryResponse["key"].(string)
		output := &data.ComplexObject{Metadata: "", Value: outputMap}
		context.SetOutput(ovOutput, output)
	}

	return true, nil
}
