package fcm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type HTTPClient struct {
	srv string
	key string

	http *http.Client
}

func NewHTTPClient(srv, key string) *HTTPClient {
	return &HTTPClient{
		srv:  srv,
		key:  key,
		http: &http.Client{},
	}
}

func (c *HTTPClient) SetTimeout(t time.Duration) {
	c.http.Timeout = t
}

func (c *HTTPClient) SendJSONRaw(message []byte) (*Response, error) {
	request, err := http.NewRequest("POST", c.srv, bytes.NewBuffer(message))

	request.Header.Set("Authorization", "key="+c.key)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return nil, wrapError(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, wrapError(err)
	}

	fcmResponse := &Response{
		RawResponse: body,
	}

	if head := response.Header.Get("Retry-After"); len(head) > 0 {
		if retry, err := strconv.Atoi(head); err == nil {
			fcmResponse.RetryAfter = time.Duration(retry) * time.Second
		} else if retry, err := http.ParseTime(head); err == nil {
			fcmResponse.RetryAfter = retry.Sub(time.Now())
		}
	}

	switch response.StatusCode {
	case http.StatusUnauthorized:
		fcmResponse.Status = RESPONSE_STATUS_UNAUTHORIZED
	case http.StatusBadRequest:
		fcmResponse.Status = RESPONSE_STATUS_BAD_REQUEST
	case http.StatusOK:
		fcmResponse.Status = RESPONSE_STATUS_OK

		if err = json.Unmarshal(body, fcmResponse); err != nil {
			return fcmResponse, wrapError(err)
		}
	}

	return fcmResponse, nil
}
