package jwtin

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/elevenia/product-simple/internal/error"
	"github.com/elevenia/product-simple/internal/options"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type ClaimsContext struct {
	echo.Context
	Claims Claims
}

type Claims struct {
	jwt.StandardClaims
	Token string
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type JWT interface {
	// Parse accepts string token and
	// return claimable context with error
	Parse(string) (Claims, errin.Error)
}

func verify(secretFile string) *rsa.PublicKey {
	verifyBytes, err := ioutil.ReadFile(secretFile)
	if err != nil {
		panic(err)
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}
	return verifyKey
}

// InitAuth auth initialization
func InitAuth(f optin.WithFlag, secretFile string) JWT {
	if !f.UseToken {
		return &jwtAuth{debug: f.Debug}
	}
	return &jwtAuth{
		publicKey: verify(secretFile),
		debug:     f.Debug,
	}
}

type jwtAuth struct {
	debug     bool
	publicKey *rsa.PublicKey
}

func (m *jwtAuth) Parse(token string) (Claims, errin.Error) {
	var claims Claims

	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return claims, errin.NewError(http.StatusForbidden, errin.ErrMalformedToken)
	}

	parsedToken, err := jwt.ParseWithClaims(
		parts[1], &claims, func(*jwt.Token) (interface{}, error) {
			return m.publicKey, nil
		},
	)
	if err != nil {
		if !m.debug {
			err = errin.ErrInvalidToken
		}
		return claims, errin.NewError(http.StatusForbidden, err)
	}
	if !parsedToken.Valid {
		return claims, errin.NewError(http.StatusForbidden, errin.ErrInvalidToken)
	}

	claims.Token = parsedToken.Raw

	return claims, nil
}
