package queryticket

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

var connectionJSON = []byte(`{
	"id" : "JiraTestConnection",
	"name": "jiraconnection",
	"description" : "JIRA Test Connection",
	"title": "JIRA Connector",
	"type": "flogo:connector",
	"version": "1.0.0",
	"ref": "Jira/connector/jiraconnection",
	"keyfield": "name",
	"settings": [
		{
		  "name": "name",
		  "value": "MyTestConnection",
		  "type": "string"
		},
		{
		  "name": "description",
		  "value": "My JIRA test Connection",
		  "type": "string"
		},
		{
		  "name": "domain",
		  "value": "https://devjira.tibco.com",
		  "type": "string"
		  
		},
		{
		  "name": "userName",
		  "value": "eseconnectors",
		  "type": "string"
		  
		},  
		{
		  "name": "password",
		  "value": "3s3c0nnectors",
		  "type": "string"
		  
		}
	  ]
}`)

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	conn := make(map[string]interface{})
	err := json.Unmarshal([]byte(connectionJSON), &conn)
	assert.Nil(t, err)

	tc.SetInput(ivConnection, conn)
	//tc.SetInput(ivQueryBy, "Recently Updated")
	tc.SetInput(ivQueryBy, "Recently Created")
	//tc.SetInput(ivProject, "ESEC")
	tc.SetInput(ivIssueType, "Enhancement")
	tc.SetInput(ivWithinTime, "120d")

	complex := &data.ComplexObject{Metadata: `{"type":"object","properties":{"project":{"required":"false", "type":"string"}}}`, Value: `{"project":"ESEC"}`}
	tc.SetInput(ivQueryParams, complex)
	act.Eval(tc)

	//name := tc.GetOutput(name).(string)

	//if name != "Test" {
	//	t.Error("Name did not match")
	//}
}
