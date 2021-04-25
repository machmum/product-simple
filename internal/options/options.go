package optin

import (
	"time"

	"github.com/labstack/echo/middleware"
)

type Options struct {
	ServerAddr  string
	Timeout     time.Duration
	BasicAuthFn middleware.BasicAuthValidator
	Flag        WithFlag
	Credential  WithCredential
}

type WithFlag struct {
	Debug        bool
	UseToken     bool
	UseBasicAuth bool
}

type WithCredential struct {
	JwtToken string
}
