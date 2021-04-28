package srvin

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/elevenia/product-simple/internal/middleware"
	"github.com/elevenia/product-simple/internal/options"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server interface {
	AddMiddleware(...echo.MiddlewareFunc)
	Echo() *echo.Echo
	Start()
}

// New returns new server with default middleware
func New(opt optin.Options) Server {
	mw := mwin.NewMiddleware(opt, opt.Credential.JwtToken)

	e := echo.New()
	e.HideBanner = true
	e.Use(
		// Recover from panics,
		// see: https://echo.labstack.com/middleware/recover
		middleware.Recover(),

		// Request ID middleware generates a unique id for a request
		middleware.RequestID(),

		// CORS handler
		mw.CORS,

		// AccessLog print information about each HTTP request
		mw.AccessLog(),
	)

	// BasicAuth implement basic auth
	if opt.Flag.UseBasicAuth {
		e.Use(middleware.BasicAuthWithConfig(mw.BasicAuthConfig()))
	}

	// ClaimsToken validate and claims object in token
	if opt.Flag.UseToken {
		e.Use(mw.AuthToken())
	}

	return &server{e: e, opt: opt}
}

// NewWithoutMiddleware return new server without middleware
func NewWithoutMiddleware(opt optin.Options) Server {
	e := echo.New()
	e.HideBanner = true
	return &server{e: e, opt: opt}
}

type server struct {
	e   *echo.Echo
	opt optin.Options
}

func (s *server) AddMiddleware(mw ...echo.MiddlewareFunc) {
	s.e.Use(mw...)
}

func (s *server) Echo() *echo.Echo {
	return s.e
}

func (s *server) Start() {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", s.opt.ServerAddr),
	}

	// Start server
	go func() {
		err := s.e.StartServer(srv)
		if err != nil {
			fmt.Printf("[server] Failed to start: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	deadline := 10 * time.Second
	fmt.Println("[server] Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	err := s.e.Shutdown(ctx)
	if err != nil {
		s.e.Logger.Fatal(err)
	}
}
