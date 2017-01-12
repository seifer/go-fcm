package fcm_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	fcm "github.com/seifer/go-fcm"
)

type config struct {
	To  string `json:"to"`
	Key string `json:"key"`
	Srv string `json:"srv"`
}

var cfg = &config{}

func init() {
	if _, err := os.Stat("test.json"); os.IsNotExist(err) {
		fmt.Println(`
            Please, you need to create file test.json for running tests. File should looks like
            {
                "to": "[TEST APP FCM REGISTRATION TOKEN]",
                "key": "[YOU FCM SERVER KEY SHOULD BE HERE]",
                "srv": "https://fcm.googleapis.com/fcm/send"
            }
        `)
		os.Exit(1)
	}

	file, err := ioutil.ReadFile("test.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		panic(err)
	}
}

func TestIncorrectDial(t *testing.T) {
	client := fcm.NewHTTPClient("stub", "blabla")

	response, err := client.SendJSONRaw([]byte("null"))

	AssertNotNil(t, err, err.Error())
	AssertIsNil(t, response)
}

func TestIncorrectAuth(t *testing.T) {
	client := fcm.NewHTTPClient(cfg.Srv, "blabla")

	response, err := client.SendJSONRaw([]byte("null"))

	AssertIsNil(t, err)
	AssertNotNil(t, response)
	Assert(t, response.Status, fcm.RESPONSE_STATUS_UNAUTHORIZED)
}

func TestIncorrectJSON(t *testing.T) {
	client := fcm.NewHTTPClient(cfg.Srv, cfg.Key)

	response, err := client.SendJSONRaw([]byte("null"))

	AssertIsNil(t, err)
	Assert(t, response.Status, fcm.RESPONSE_STATUS_BAD_REQUEST)
}

func TestCorrectRequestNotRegistered(t *testing.T) {
	client := fcm.NewHTTPClient(cfg.Srv, cfg.Key)

	response, err := client.SendJSONRaw(getTestMessage(cfg.To))

	AssertIsNil(t, err)
	Assert(t, response.Status, fcm.RESPONSE_STATUS_OK, fmt.Sprintf("Status %s", response.RawResponse))
	Assert(t, response.Success, 1, fmt.Sprintf("Success %s", response.RawResponse))
	Assert(t, response.Failure, 0, fmt.Sprintf("Failure %s", response.RawResponse))
	Assert(t, response.CanonicalIds, 0, fmt.Sprintf("CanonicalIds %s", response.RawResponse))
	Assert(t, len(response.Results), 1, fmt.Sprintf("Results %s", response.RawResponse))
}

func TestCorrectRequestUserRegistered(t *testing.T) {
	client := fcm.NewHTTPClient(cfg.Srv, cfg.Key)

	response, err := client.SendJSONRaw(getTestMessage(""))

	AssertIsNil(t, err)
	Assert(t, response.Status, fcm.RESPONSE_STATUS_OK, fmt.Sprintf("Status %s", response.RawResponse))
	Assert(t, response.Success, 1, fmt.Sprintf("Success %s", response.RawResponse))
	Assert(t, response.Failure, 0, fmt.Sprintf("Failure %s", response.RawResponse))
	Assert(t, response.CanonicalIds, 0, fmt.Sprintf("CanonicalIds %s", response.RawResponse))
	Assert(t, len(response.Results), 1, fmt.Sprintf("Results %s", response.RawResponse))
}
