package fcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	ERROR_RESPONSE_CODE_MISSING_REGISTRATION         = "MissingRegistration"
	ERROR_RESPONSE_CODE_INVALID_REGISTRATION         = "InvalidRegistration"
	ERROR_RESPONSE_CODE_NOT_REGISTERED               = "NotRegistered"
	ERROR_RESPONSE_CODE_INVALID_PACKAGE_NAME         = "InvalidPackageName"
	ERROR_RESPONSE_CODE_MISMATCH_SENDERID            = "MismatchSenderId"
	ERROR_RESPONSE_CODE_MESSAGE_TOO_BIG              = "MessageTooBig"
	ERROR_RESPONSE_CODE_INVALID_DATA_KEY             = "InvalidDataKey"
	ERROR_RESPONSE_CODE_INVALID_TTL                  = "InvalidTtl"
	ERROR_RESPONSE_CODE_TIMEOUT                      = "Unavailable"
	ERROR_RESPONSE_CODE_INTERNAL_SERVER_ERROR        = "InternalServerError"
	ERROR_RESPONSE_CODE_DEVICE_MESSAGE_RATE_EXCEEDED = "DeviceMessageRateExceeded"
	ERROR_RESPONSE_CODE_TOPICS_MESSAGE_RATE_EXCEEDED = "TopicsMessageRateExceeded"
)

type result struct {
	MessageId      string `json:"message_id"`
	RegistrationId string `json:"registration_id"`
	Error          string `json:"error"`
}

type HTTPClient struct {
	srv string
	key string
}

type HTTPResponse struct {
	// HTTP data
	RetryAfter      time.Duration
	HTTPStatusCode  int
	HTTPRawResponse []byte

	// Original body response
	MulticastId  int      `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIds int      `json:"canonical_ids"`
	Results      []result `json:"results"`
}

func NewHTTPClient(srv, key string) *HTTPClient {
	return &HTTPClient{srv, key}
}

func (c *HTTPClient) SendJSONRaw(message []byte) (*HTTPResponse, error) {
	request, err := http.NewRequest("POST", c.srv, bytes.NewBuffer(message))

	request.Header.Set("Authorization", "key="+c.key)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusBadRequest {
		return nil, errors.New(
			fmt.Sprintf("HTTP request failed. Status %d. %s", response.StatusCode, body),
		)
	}

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.New(
			fmt.Sprintf("HTTP request failed. Status %d. Authentication Error", response.StatusCode),
		)
	}

	fcmResponse := &HTTPResponse{
		HTTPStatusCode:  response.StatusCode,
		HTTPRawResponse: body,
	}

	if head := response.Header.Get("Retry-After"); len(head) > 0 {
		if retry, err := strconv.Atoi(head); err == nil {
			fcmResponse.RetryAfter = time.Duration(retry) * time.Second
		} else if retry, err := http.ParseTime(head); err == nil {
			fcmResponse.RetryAfter = retry.Sub(time.Now())
		}
	}

	if response.StatusCode == http.StatusOK {
		if err = json.Unmarshal(body, fcmResponse); err != nil {
			return fcmResponse, err
		}
	}

	return fcmResponse, nil
}
