package analysis

import (
	"github.com/labstack/echo/v4"

	"github.com/edersonbrilhante/ccvs"
)

// Create creates a new Analysis
func (a Analysis) Create(c echo.Context, req ccvs.Analysis) (ccvs.Analysis, error) {
	return a.adb.Create(a.db, req)
}

// View returns single Analysis
func (a Analysis) View(c echo.Context, id string) (ccvs.Analysis, error) {
	return a.adb.View(a.db, id)
}
