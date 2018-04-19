package add

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLog is the default logger for the Add Numbers Activity
var activityLog = logger.GetLogger("activity-Asha-addNumbers")

const (
	ivNum1 = "Number1"
	ivNum2 = "Number2"

	ovAddition = "Addition"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

// LogActivity is an Activity that is used to log a message to the console
// inputs : {message, flowInfo}
// outputs: none
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

// Eval implements api.Activity.Eval - Logs the Message
func (a *AddActivity) Eval(context activity.Context) (done bool, err error) {

	//mv := context.GetInput(ivMessage)
	num1, _ := context.GetInput(ivNum1).(int)
	num2, _ := context.GetInput(ivNum2).(int)

	activityLog.Info(fmt.Sprintf("Num1: %d, Num2: %d", num1, num2))
	activityLog.Info(fmt.Sprintf("Addition is : %d", num1+num2))
	context.SetOutput(ovAddition, num1+num2)

	return true, nil
}
