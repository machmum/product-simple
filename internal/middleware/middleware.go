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
)

func NewMiddleware(f optin.WithFlag, secretFile string) *mw {
	return &mw{b: jwtin.InitAuth(f, secretFile)}
}

// mw middleware contains jwt pub key
type mw struct {
	b jwtin.JWT
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

func (m *mw) AuthToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			token := c.Request().Header.Get(constant.HeaderRequestedToken)
			if token == "" {
				logger.InfoLog(context.Background(), constant.Failed, "AuthToken.GetToken", errin.ErrTokenRequired)
			}

			claims, err := m.b.Parse(token)
			if err != nil {
				logger.InfoLog(context.Background(), constant.Failed, "AuthToken.Parse", err)
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
