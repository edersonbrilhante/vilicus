package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewHTTP creates new analysis http service
func NewHTTP(r *echo.Group) {
	h := HTTP{}
	ur := r.Group("/healthz")

	ur.GET("", h.view)
}

// HTTP represents analysis http service
type HTTP struct{}

func (h HTTP) view(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING\n")
}
