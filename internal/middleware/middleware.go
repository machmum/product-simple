package mwin

import (
	"context"
	"net/http"
	"strings"

	"github.com/elevenia/product-simple/internal/constant"
	"github.com/elevenia/product-simple/internal/error"
	"github.com/elevenia/product-simple/internal/log"
	"github.com/elevenia/product-simple/internal/options"
	"github.com/elevenia/product-simple/internal/token/jwt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// failed will write a default template response when returning a failed response
func failed(c echo.Context, code int, err error) error {
	response := map[string]interface{}{
		"data": nil,
		"error": map[string]interface{}{
			"code":    0,
			"status":  code,
			"message": err.Error(),
		},
	}
	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(code, response)
}

func hasAuthorization(c echo.Context) bool {
	return c.Request().Header.Get(constant.HeaderAuth) != ""
}

func isBearerAuth(c echo.Context) bool {
	if !hasAuthorization(c) {
		return false
	}
	parts := strings.SplitN(c.Request().Header.Get(constant.HeaderAuth), " ", 2)
	if !(len(parts) == 2 && parts[0] == constant.AuthBearer) {
		return false
	}
	return true
}

func NewMiddleware(opt optin.Options, secretFile string) *mw {
	return &mw{b: jwtin.InitAuth(opt.Flag, secretFile), opt: opt}
}

// mw middleware contains jwt pub key
type mw struct {
	b             jwtin.JWT
	opt           optin.Options
	basicAuthSkip bool
}

func (m *mw) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(constant.HeaderAllowOrigin, "*")
		c.Response().Header().Set(constant.HeaderAllowMethod, strings.Join(constant.AllowedMethods, ","))
		c.Response().Header().Set(constant.HeaderAllowHeaders, strings.Join(constant.AllowedHeaders, ","))

		if c.Request().Method == http.MethodOptions {
			_, _ = c.Response().Write([]byte("ok"))
		}

		return next(c)
	}
}

func (m *mw) skipBasicAuth() bool {
	m.basicAuthSkip = true
	return true
}

func (m *mw) useBasicAUth() bool {
	m.basicAuthSkip = false
	return false
}

func (m *mw) BasicAuthConfig() middleware.BasicAuthConfig {
	skipper := func(c echo.Context) bool {
		if m.opt.Flag.UseBasicAuth {
			if m.opt.Flag.UseToken && isBearerAuth(c) {
				return m.skipBasicAuth()
			}
			return m.useBasicAUth()
		}
		return m.skipBasicAuth()
	}

	return middleware.BasicAuthConfig{
		Skipper:   skipper,
		Validator: m.opt.BasicAuthFn,
	}
}

func (m *mw) AuthToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if m.opt.Flag.UseToken && !isBearerAuth(c) {
				if m.opt.Flag.UseBasicAuth && !m.basicAuthSkip {
					return next(&jwtin.ClaimsContext{Context: c})
				}
				return failed(c, http.StatusUnauthorized, errin.ErrTokenRequired)
			}

			token := c.Request().Header.Get(constant.HeaderAuth)
			if token == "" {
				logger.InfoLog(context.Background(), constant.Failed, "AuthToken.GetToken", errin.ErrTokenRequired)
				return failed(c, http.StatusUnauthorized, errin.ErrTokenRequired)
			}

			claims, err := m.b.Parse(token)
			if err != nil {
				logger.InfoLog(context.Background(), constant.Failed, "AuthToken.Parse", err)
				return failed(c, err.Code(), err)
			}

			return next(&jwtin.ClaimsContext{
				Context: c,
				Claims:  claims,
			})
		}
	}
}

// AccessLog: log session's email & basic request properties data
func (m *mw) AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)

			if c.Request().Header.Get(constant.HeaderRequestID) != "" {
				reqID = c.Request().Header.Get(constant.HeaderRequestID)
			}

			// get email from ClaimsContext
			// embed as `_email` in log
			var email string
			if get, ok := c.(*jwtin.ClaimsContext); ok {
				email = get.Claims.Email
			}

			defer func() {
				rq := c.Request()
				rs := c.Response()
				ss := map[string]interface{}{"_email": email}

				logger.AccessLog(reqID, ss, rq, rs).Info(logger.LLvlAccess)
			}()
			return next(c)
		}
	}
}
