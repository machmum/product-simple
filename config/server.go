package config

import (
	cfg "github.com/elevenia/product-simple/config/env"
	"github.com/elevenia/product-simple/internal/options"
	"github.com/labstack/echo"
)

func ServerOpts() optin.Options {
	return optin.Options{
		ServerAddr: cfg.Env.ServerAddress,
		BasicAuthFn: func(username, password string, c echo.Context) (valid bool, err error) {
			if username == cfg.Env.BasicAuthUsername &&
				password == cfg.Env.BasicAuthUsername {
				valid = true
			}
			return
		},
		Flag: optin.WithFlag{
			Debug:        cfg.Env.Debug,
			UseToken:     false,
			UseBasicAuth: false,
		},
		Credential: optin.WithCredential{
			JwtToken: cfg.Env.JwtPublicKey,
		},
	}
}
