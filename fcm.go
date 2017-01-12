package fcm

import (
	"fmt"
	"time"
)

const (
	// status of response
	RESPONSE_STATUS_UNAUTHORIZED = 1
	RESPONSE_STATUS_BAD_REQUEST  = 2
	RESPONSE_STATUS_OK           = 3
	RESPONSE_STATUS_SERVER_ERROR = 4

	// 200 OK with error at message
	RESPONSE_MESSAGE_CODE_MISSING_REGISTRATION         = "MissingRegistration"
	RESPONSE_MESSAGE_CODE_INVALID_REGISTRATION         = "InvalidRegistration"
	RESPONSE_MESSAGE_CODE_NOT_REGISTERED               = "NotRegistered"
	RESPONSE_MESSAGE_CODE_INVALID_PACKAGE_NAME         = "InvalidPackageName"
	RESPONSE_MESSAGE_CODE_MISMATCH_SENDERID            = "MismatchSenderId"
	RESPONSE_MESSAGE_CODE_MESSAGE_TOO_BIG              = "MessageTooBig"
	RESPONSE_MESSAGE_CODE_INVALID_DATA_KEY             = "InvalidDataKey"
	RESPONSE_MESSAGE_CODE_INVALID_TTL                  = "InvalidTtl"
	RESPONSE_MESSAGE_CODE_TIMEOUT                      = "Unavailable"
	RESPONSE_MESSAGE_CODE_INTERNAL_SERVER_ERROR        = "InternalServerError"
	RESPONSE_MESSAGE_CODE_DEVICE_MESSAGE_RATE_EXCEEDED = "DeviceMessageRateExceeded"
	RESPONSE_MESSAGE_CODE_TOPICS_MESSAGE_RATE_EXCEEDED = "TopicsMessageRateExceeded"
)

type result struct {
	MessageId      string `json:"message_id"`
	RegistrationId string `json:"registration_id"`
	Error          string `json:"error"`
}

type Response struct {
	// Custom data
	Status int

	// HTTP data
	RetryAfter  time.Duration
	RawResponse []byte

	// Original body response
	MulticastId  int      `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIds int      `json:"canonical_ids"`
	Results      []result `json:"results"`
}

func (r *Response) String() string {
	return fmt.Sprintf("{Status: %d RetryAfter: %s Body: %s", r.Status, r.RetryAfter, string(r.RawResponse))
}
