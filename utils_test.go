package fcm_test

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func Assert(t *testing.T, obtained, expected interface{}, comment ...string) {
	assert(t, obtained, expected, comment...)
}

func AssertIsNil(t *testing.T, obtained interface{}, comment ...string) {
	if obtained == nil {
		return
	}

	assert(t, reflect.ValueOf(obtained).IsNil(), true, comment...)
}

func AssertNotNil(t *testing.T, obtained interface{}, comment ...string) {
	if obtained == nil {
		assert(t, true, false, comment...)
		return
	}

	assert(t, reflect.ValueOf(obtained).IsNil(), false, comment...)
}

// Utils
func assert(t *testing.T, obtained, expected interface{}, comment ...string) {
	if reflect.DeepEqual(expected, obtained) {
		return
	}

	title := "UNKNOWN TEST"

	if pc, _, _, ok := runtime.Caller(2); ok {
		title = getFuncName(pc)
	}

	if len(comment) == 0 {
		t.Errorf("\n%s\nExpected=%v\nObtained=%v\n", title, expected, obtained)
	} else {
		t.Errorf("\n%s\nExpected=%v\nObtained=%v\n%s\n", title, expected, obtained, comment[0])
	}

	t.FailNow()
}

func getFuncName(pc uintptr) string {
	f := runtime.FuncForPC(pc)
	file, line := f.FileLine(pc)

	full := f.Name()

	for i := len(full) - 1; i > 0; i-- {
		if full[i] == '.' {
			full = string(full[i+1:])
			break
		}
	}

	full += " in " + file + " line " + fmt.Sprintf("%d", line)

	return full
}

func getTestMessage(to string) []byte {
	if to == "" {
		to = cfg.To
	}

	return []byte(`{
		"to": "` + cfg.To + `",
		"dry_run": true,
		"priority": "high",
		"time_to_live": 60,
		"notification": {
			"title": "Test title"
		}
    }`)
}
