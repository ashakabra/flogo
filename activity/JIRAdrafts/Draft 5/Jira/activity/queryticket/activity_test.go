package queryticket

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

var activityMetadata *activity.Metadata

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
	tc.SetInput(ivDomain, "https://jira.tibco.com")
	tc.SetInput(ivBasicAuthToken, "YWthYnJhOkxpZ2h0MTIjJA==")
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
