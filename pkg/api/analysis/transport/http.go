package transport

import (	
	"github.com/edersonbrilhante/ccvs"
	"github.com/edersonbrilhante/ccvs/pkg/api/analysis"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Analysis create request
type createReq struct {
	Image string `json:"image" validate:"required"`
}

// NewHTTP creates new analysis http service
func NewHTTP(svc analysis.Service, r *echo.Group) {
	h := HTTP{svc}
	ur := r.Group("/analysis")

	ur.POST("", h.create)
	ur.GET("/:id", h.view)
}

// HTTP represents analysis http service
type HTTP struct {
	svc analysis.Service
}

func (h HTTP) create(c echo.Context) error {
	req := new(createReq)

	if err := c.Bind(req); err != nil {
		return err
	}

	analysis, err := h.svc.Create(c, ccvs.Analysis{
		Image: req.Image,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, analysis)
}

func (h HTTP) view(c echo.Context) error {
	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	return gorsk.ErrBadRequest
	// }

	// result, err := h.svc.View(c, id)
	// if err != nil {
	// 	return err
	// }
	result := "view"

	return c.JSON(http.StatusOK, result)
}
