package queryticket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var activityLog = logger.GetLogger("jira-activity-queryticket")

const (
	queryByUpdate = "Recently Updated"
	queryByCreate = "Recently Created"

	ivDomain         = "domain"
	ivBasicAuthToken = "basicAuthToken"
	ivQueryBy        = "queryBy"
	ivProject        = "project"
	ivIssueType      = "issueType"
	ivWithinTime     = "withinTime"
	ivQueryParams    = "queryParams"

	ovOutput = "output"
)

type QueryTicketActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &QueryTicketActivity{metadata: metadata}
}

func (a *QueryTicketActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func ParseOutput(outputSchema map[string]interface{}) ([]string, error) {
	if outputSchema == nil {
		return nil, nil
	}

	props := outputSchema["items"].(map[string]interface{})
	properties := props["properties"].(map[string]interface{})
	fields := make([]string, len(properties))
	i := 0
	for k, _ := range properties {
		fields[i] = k
		i++
	}
	return fields, nil
}

func (a *QueryTicketActivity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Info("JIRA Query Ticket")
	domain := context.GetInput(ivDomain).(string)
	basicAuthToken := context.GetInput(ivBasicAuthToken).(string)
	queryBy := context.GetInput(ivQueryBy).(string)
	project := context.GetInput(ivProject).(string)
	issueType := context.GetInput(ivIssueType).(string)
	withinTime := context.GetInput(ivWithinTime).(string)
	parameters, err := GetParameter(context.GetInput(ivQueryParams))

	outputMap, _ := LoadJsonSchemaFromMetadata(context.GetOutput(ovOutput))
	if outputMap != nil {
		outputFields, _ := ParseOutput(outputMap)
		activityLog.Infof("Reading Output is :: %s", outputFields)
	}

	fmt.Printf("Input Values are %s, %s, %s, %s, %s, %s", domain, basicAuthToken, queryBy, project, issueType, withinTime)

	var input string
	if queryBy == queryByUpdate {
		//input = "project='" + project + "' AND issueType='" + issueType + "' AND updated >= -" + withinTime
		input = "issueType='" + issueType + "' AND updated >= -" + withinTime
	} else if queryBy == queryByCreate {
		//input = "project='" + project + "' AND issueType='" + issueType + "' AND created >= -" + withinTime
		input = "issueType='" + issueType + "' AND created >= -" + withinTime
	}

	for _, value := range parameters.QueryParams {
		activityLog.Infof("Name :: %s , Value :: %s ", value.Name, value.Value)
		stringValue := fmt.Sprint(value.Value)
		input = input + " AND " + value.Name + " in (" + stringValue + ")"
	}

	activityLog.Infof("Input Value :: %s", input)
	url := domain + "/rest/api/2/search?jql=" + url.QueryEscape(input)

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Basic "+basicAuthToken)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		jsonResponseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return false, fmt.Errorf("Error reading JSON response data after query invocation, %s", err.Error())
		}
		defer response.Body.Close()

		var queryResponse interface{}
		err = json.Unmarshal(jsonResponseData, &queryResponse)

		if err != nil {
			fmt.Printf("Error :: %s", err)
		}

		m := queryResponse.(map[string]interface{})
		issues := m["issues"].([]interface{})
		responseIssues := make([]map[string]interface{}, len(issues))
		if len(issues) > 0 {
			for i := range issues {
				issue := issues[i].(map[string]interface{})
				responseIssues[i] = make(map[string]interface{})
				responseIssues[i]["key"] = issue["key"].(string)
				responseIssues[i]["summary"] = issue["fields"].(map[string]interface{})["summary"]
			}
		} else {
			activityLog.Infof("No issues found")
		}

		//activityLog.Infof("Output is -- %s", responseIssues)
		output := &data.ComplexObject{Metadata: "", Value: responseIssues}
		context.SetOutput(ovOutput, output)
	}

	return true, nil
}
