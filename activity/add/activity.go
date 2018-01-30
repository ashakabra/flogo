package add

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-aasha-add")

const (
	ivNum1   = "Number1"
	ivNum2  = "Number2"
	
	ovSum = "Addition"
)

// AddActivity is an Activity that is used to add two numbers
// inputs : {number1, number2}
// outputs: output
type AddActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AddActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *AddActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Adds the 2 input numbers
func (a *AddActivity) Eval(context activity.Context) (done bool, err error) {

	//mv := context.GetInput(ivMessage)
	num1, _ := context.GetInput(ivNum1).(int)
	num2, _ := context.GetInput(ivNum2).(int)

	//activityLog.Info("Number one is::" +num1+" and Number two is ::"+num2)

	context.SetOutput(ovSum, num1+num2)
	
	return true, nil
}

