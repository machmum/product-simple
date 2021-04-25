package response

import (
	"fmt"
	"net/http"

	"github.com/elevenia/product-simple/internal/constant"
	"github.com/labstack/echo"
)

var debug bool

func Init(d bool) {
	debug = d
}

func mappingError(resCode int, err error) (int, string) {
	if debug && err != nil {
		return resCode, err.Error()

	} else {
		switch resCode {
		case http.StatusOK:
			return resCode, ""
		case http.StatusBadRequest:
			return resCode, "Missing request"
		case http.StatusUnauthorized:
			return resCode, "User is not authorized"
		case http.StatusForbidden:
			return resCode, "Access Denied"
		case http.StatusNotFound:
			return resCode, "Data Not Found"
		case http.StatusMethodNotAllowed:
			return resCode, "Action Not Allowed"
		case http.StatusConflict:
			return resCode, "Data Conflict"
		case http.StatusUnprocessableEntity:
			return resCode, "Data can't be processed"
		case http.StatusInternalServerError:
			return resCode, "Internal server error"
		case http.StatusInsufficientStorage:
			return resCode, "Error Database"
		default:
			return resCode, fmt.Sprintf("Internal service error")
		}
	}
}

// response struct
type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Ref     string      `json:"ref,omitempty"`
}

func JSON(c echo.Context, code int, result interface{}, err error) error {
	resp := new(Response)
	if err != nil {
		reqID := c.Response().Header().Get(echo.HeaderXRequestID)
		if c.Request().Header.Get(constant.HeaderRequestID) != "" {
			reqID = c.Request().Header.Get(constant.HeaderRequestID)
		}
		resp.Ref = reqID
	}

	resp.Code, resp.Message = mappingError(code, err)
	resp.Data = result

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(code, resp)
}
