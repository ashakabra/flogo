package createticket

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	createRestAPI = "/rest/api/2/issue/"
)

//Connection datastructure for storing jira connection details
type Connection struct {
	Domain   string
	UserName string
	Password string
}

//GetConnection returns a deserialized connection object
func GetConnection(connector interface{}) *Connection {
	connectionInfo := connector.(map[string]interface{})
	connectionSettings := connectionInfo["settings"].([]interface{})
	conn := &Connection{}
	for _, v := range connectionSettings {
		setting := v.(map[string]interface{})
		switch setting["name"] {
		case "domain":
			conn.Domain = setting["value"].(string)
		case "userName":
			conn.UserName = setting["value"].(string)
		case "password":
			conn.Password = setting["value"].(string)
		}
	}
	return conn
}

//CreateTicket calls jira rest api to create ticket on jira
func CreateTicket(conn *Connection, inputJSON map[string]interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(inputJSON)
	if err != nil {
		return nil, fmt.Errorf("Cannot deserialize inputData: %s", err.Error())
	}
	activityLog.Debugf("Activity Input :: %s", jsonData)
	url := conn.Domain + createRestAPI

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Basic "+BasicAuth(conn.UserName, conn.Password))

	client := &http.Client{}

	return client.Do(request)
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
