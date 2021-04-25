package http

import (
	"context"
	"net/http"

	"github.com/elevenia/product-simple/domain/product"
	"github.com/elevenia/product-simple/helper/response"
	"github.com/elevenia/product-simple/internal/constant"
	"github.com/elevenia/product-simple/internal/context"
	"github.com/labstack/echo"
)

// Handler: represent the http handler
type Handler struct {
	uc product.Usecase
}

func NewHandler(e *echo.Echo, us product.Usecase) {
	handler := &Handler{uc: us}

	e.GET("/product", handler.view)
}

func newCtx(c echo.Context) context.Context {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	if c.Request().Header.Get(constant.HeaderRequestID) != "" {
		reqID = c.Request().Header.Get(constant.HeaderRequestID)
	}
	return ctxin.NewRequestIDContext(c.Request().Context(), reqID)
}

func (h *Handler) view(c echo.Context) error {
	var ctx = newCtx(c)

	result, err := h.uc.View(ctx)
	if err != nil {
		return response.JSON(c, err.Code(), nil, err)
	}

	return response.JSON(c, http.StatusOK, result, nil)
}
