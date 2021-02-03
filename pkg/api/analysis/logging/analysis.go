package analysis

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/edersonbrilhante/vilicus/pkg/api/analysis"
	"github.com/edersonbrilhante/vilicus/pkg/log"
	"github.com/edersonbrilhante/vilicus/pkg/types"
)

// New creates new analysis logging service
func New(svc analysis.Service, logger log.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents analysis logging service
type LogService struct {
	analysis.Service
	logger log.Logger
}

const name = "analysis"

// Create logging
func (ls *LogService) Create(c echo.Context, req types.Analysis) (resp types.Analysis, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create analysis request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req types.Analysis) (resp types.Analysis, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Update analysis request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, req)
}
