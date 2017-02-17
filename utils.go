package fcm

import (
	"fmt"
	"runtime"
)

type Error int

func (e *Error) Error() string {
	return ""
}

func wrapError(err error) error {
	var funcName string

	if pc, _, _, ok := runtime.Caller(1); ok {
		funcName = getFuncName(pc)
	}

	return fmt.Errorf("FCM %s error: %s", funcName, err)
}
