package constant

import (
	"net/http"
)

type Status string

const (
	Received  Status = "received"
	Published Status = "published"
	Sent      Status = "sent"
	Failed    Status = "failed"
	Error     Status = "error"
	Success   Status = "success"
)

const (
	HeaderContentType    = "Content-Type"
	HeaderRequestID      = "Access-Control-Requested-Id"
	HeaderAllowOrigin    = "Access-Control-Allow-Origin"
	HeaderAllowMethod    = "Access-Control-Allow-Methods"
	HeaderAllowHeaders   = "Access-Control-Allow-Headers"
	HeaderAuth           = "Authorization"
	HeaderPermission     = "Access-Control-Requested-For"
	HeaderRequestedToken = "Access-Control-Requested-Token"
	HeaderRequestedHost  = "Access-Control-Requested-Host"

	// time template
	DatetimeYMDFormat = "2006-01-02 15:04:05"
	DatetimeDMYFormat = "02-01-2006 15:04:05"
)

var AllowedHeaders = []string{
	HeaderContentType,
	HeaderPermission,
	HeaderRequestID,
	HeaderAllowOrigin,
	HeaderAllowMethod,
	HeaderAllowHeaders,
	HeaderAuth,
	HeaderRequestedToken,
	HeaderRequestedHost,
}

var AllowedMethods = []string{
	http.MethodPost,
	http.MethodGet,
	http.MethodOptions,
	http.MethodPut,
	http.MethodDelete,
	http.MethodPatch,
}
