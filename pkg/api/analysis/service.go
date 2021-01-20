package analysis

import (
	"github.com/edersonbrilhante/ccvs"
	"github.com/edersonbrilhante/ccvs/pkg/api/analysis/platform/pgsql"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/labstack/echo/v4"
)

// Service represents analysis application interface
type Service interface {
	Create(echo.Context, ccvs.Analysis) (ccvs.Analysis, error)
	View(echo.Context, string) (ccvs.Analysis, error)
}

// New creates new analysis application service
func New(db *pg.DB, adb ADB) *Analysis {
	return &Analysis{db: db, adb: adb}
}

// Initialize initalizes Analysis application service with defaults
func Initialize(db *pg.DB) *Analysis {
	return New(db, pgsql.Analysis{})
}

// Analysis represents analysis application service
type Analysis struct {
	db  *pg.DB
	adb ADB
}

// // Create creates a new analysis
// func (a Analysis) Create(c echo.Context, req ccvs.Analysis) (ccvs.Analysis, error) {
// 	// if err := nil; err != nil {
// 	// 	return ccvs.Analysis{}, err
// 	// }
// 	return req, nil
// }

// ADB represents analysis repository interface
type ADB interface {
	Create(orm.DB, ccvs.Analysis) (ccvs.Analysis, error)
	View(orm.DB, string) (ccvs.Analysis, error)
}
