package fcm_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
var stubToken = "fNom830_YZA:APA91bGWbm3rVyCv1DuD3FyqoExoeHPs_nAGm2WoJjcc-JNeUHhBaQqGbjbrku7DPBW9LTajXwCLD6rQwD2fa11dxUvwdbZzd41fwVMtZj9ZlhE3a8yWkDpWP6z3uEVh-0AUd2rkjzuE"

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

func TestIncorrectAuth(t *testing.T) {
	client := fcm.NewHTTPClient(cfg.Srv, "blabla")

	_, err := client.SendJSONRaw([]byte("null"))

	AssertNotNil(t, err, err.Error())
}

func TestIncorrectJSON(t *testing.T) {
	client := fcm.NewHTTPClient(cfg.Srv, cfg.Key)

	_, err := client.SendJSONRaw([]byte("null"))

	AssertNotNil(t, err, err.Error())
}

func TestCorrectRequestNotRegistered(t *testing.T) {
	client := fcm.NewHTTPClient(cfg.Srv, cfg.Key)

	response, err := client.SendJSONRaw(getTestMessage(stubToken))

	AssertIsNil(t, err)
	Assert(t, response.HTTPStatusCode, http.StatusOK, fmt.Sprintf("Status %s", response.HTTPRawResponse))
	Assert(t, response.Success, 1, fmt.Sprintf("Success %s", response.HTTPRawResponse))
	Assert(t, response.Failure, 0, fmt.Sprintf("Failure %s", response.HTTPRawResponse))
	Assert(t, response.CanonicalIds, 0, fmt.Sprintf("CanonicalIds %s", response.HTTPRawResponse))
	Assert(t, len(response.Results), 1, fmt.Sprintf("Results %s", response.HTTPRawResponse))
}

func TestCorrectRequestUserRegistered(t *testing.T) {
	client := fcm.NewHTTPClient(cfg.Srv, cfg.Key)

	response, err := client.SendJSONRaw(getTestMessage(""))

	AssertIsNil(t, err)
	Assert(t, response.HTTPStatusCode, http.StatusOK, fmt.Sprintf("Status %s", response.HTTPRawResponse))
	Assert(t, response.Success, 1, fmt.Sprintf("Success %s", response.HTTPRawResponse))
	Assert(t, response.Failure, 0, fmt.Sprintf("Failure %s", response.HTTPRawResponse))
	Assert(t, response.CanonicalIds, 0, fmt.Sprintf("CanonicalIds %s", response.HTTPRawResponse))
	Assert(t, len(response.Results), 1, fmt.Sprintf("Results %s", response.HTTPRawResponse))
}
