package logger

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

func getCaller(ok bool, fn string, line int) (cf string) {
	if ok {
		cf = fmt.Sprintf("%v:%d", fn, line)
	}
	return
}

// getMethod: return trimmed `.go` from file & method name from fn
func getMethod(fn, file string) string {
	// get last index from file
	file = path.Base(file)

	// get string before `.func1` from fn
	base := path.Base(fn)
	if strings.LastIndex(fn, "func1") > 0 {
		base = path.Base(fn[:strings.LastIndex(fn, "func1")-1])
	}

	return fmt.Sprintf("%s.%s", file[:strings.LastIndex(file, ".")], base[strings.LastIndex(base, ".")+1:])
}

func GetFields(jsonMessage string) interface{} {
	fields := logrus.Fields{}
	_ = json.Unmarshal([]byte(jsonMessage), &fields)
	return fields
}
