package logger

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"runtime"

	"github.com/elevenia/product-simple/internal/constant"
	ctxin "github.com/elevenia/product-simple/internal/context"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const (
	// log msg
	LLvlAccess   = "access_log"
	LLvlService  = "service_log"
	LLvlPostgres = "postgres_log"
	LLvlInfo     = "info_log"
)

type Options struct {
	// Entry only supports:
	// ErrorLevel
	// WarnLevel
	// FatalLevel
	// InfoLevel
	// DebugLevel
	Level logrus.Level
}

// New: new json-formatter log
func New(opt Options) {
	logrus.SetOutput(os.Stdout)
	if opt.Level < 1 {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(opt.Level)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: constant.DatetimeYMDFormat})

	return
}

// AccessLog logs access log
func AccessLog(id string, session interface{}, request *http.Request, response *echo.Response) *logger {
	return &logger{f: accessLog(id, session, request, response)}
}

func accessLog(id string, session interface{}, request *http.Request, response *echo.Response) logrus.Fields {
	fields := logrus.Fields{
		"_session": session,
		"id":       id,
		"host":     request.Host,
		"method":   request.Method,
		"uri":      request.URL.String(),
		"status":   response.Status,
	}
	return fields
}

// ServiceLog logs param, request and response (error if any)
func ServiceLog(ctx context.Context, param, request, response interface{}, err error) Logger {
	return &logger{f: serviceLog(ctx, param, request, response, err)}
}

func serviceLog(ctx context.Context, param, request, response interface{}, err error) logrus.Fields {
	var (
		requestID = ctxin.FromRequestIDContext(ctx)
		fields    = make(logrus.Fields, 0)
		skip      = ctxin.FromSkipContext(ctx)
	)

	if skip < 1 {
		skip = 3
	}

	pc, file, line, ok := runtime.Caller(skip)
	fn := runtime.FuncForPC(pc).Name()

	fields["caller"] = getCaller(ok, fn, line)
	fields["method"] = getMethod(fn, file)

	if requestID != "" {
		fields["id"] = requestID
	}
	if param != nil {
		rb, _ := json.Marshal(map[string]interface{}{"data": param})
		fields["param"] = GetFields(string(rb))
	}
	if request != nil {
		rb, _ := json.Marshal(map[string]interface{}{"data": request})
		fields["request"] = GetFields(string(rb))
	}
	if response != nil {
		rb, _ := json.Marshal(map[string]interface{}{"data": response})
		fields["response"] = GetFields(string(rb))
	}

	if err != nil {
		fields["error"] = err
	}

	return fields
}

// InfoLogs logs with status
func InfoLog(ctx context.Context, status constant.Status, info string, err error) Logger {
	return &logger{f: infoLog(ctx, status, info, err)}
}

func infoLog(ctx context.Context, status constant.Status, info string, err error) logrus.Fields {
	var (
		requestID = ctxin.FromRequestIDContext(ctx)
		fields    = make(logrus.Fields, 0)
		skip      = ctxin.FromSkipContext(ctx)
	)

	if skip < 1 {
		skip = 2
	}

	pc, file, line, ok := runtime.Caller(skip)
	fn := runtime.FuncForPC(pc).Name()

	fields["caller"] = getCaller(ok, fn, line)
	fields["method"] = getMethod(fn, file)

	if requestID != "" {
		fields["id"] = requestID
	}

	if status != "" {
		fields["status"] = status
	}
	if info != "" {
		fields["info"] = info
	}

	if err != nil {
		fields["error"] = err
	}

	return fields
}

func SqlLog(ctx context.Context, ql QueryHelper, query string, bind ...interface{}) Logger {
	return &logger{f: sqlLog(ctx, ql, query, bind)}
}

func sqlLog(ctx context.Context, ql QueryHelper, query string, bind ...interface{}) logrus.Fields {
	var (
		requestID = ctxin.FromRequestIDContext(ctx)
		fields    = make(logrus.Fields, 0)
		skip      = ctxin.FromSkipContext(ctx)
	)

	if skip < 1 {
		skip = 3
	}

	pc, file, line, ok := runtime.Caller(skip)
	fn := runtime.FuncForPC(pc).Name()

	fields["caller"] = getCaller(ok, fn, line)
	fields["method"] = getMethod(fn, file)

	if requestID != "" {
		fields["id"] = requestID
	}

	fields["query"] = query

	if bind != nil {
		fields["parameter"] = bind[0]
	}

	if ql.Affected > 0 {
		fields["rows_affected"] = ql.Affected
	}

	if ql.Error != nil {
		fields["error"] = ql.Error
	}

	return fields
}

type QueryHelper struct {
	Error    error `json:"-"`
	Affected int64 `json:"-"`
}
