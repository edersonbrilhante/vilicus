package api

import (
	"os"

	"github.com/edersonbrilhante/ccvs/pkg/api/analysis"
	al "github.com/edersonbrilhante/ccvs/pkg/api/analysis/logging"
	at "github.com/edersonbrilhante/ccvs/pkg/api/analysis/transport"
	"github.com/edersonbrilhante/ccvs/pkg/util/config"
	"github.com/edersonbrilhante/ccvs/pkg/util/postgres"
	"github.com/edersonbrilhante/ccvs/pkg/util/server"
	"github.com/edersonbrilhante/ccvs/pkg/util/zlog"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {

	// // db connection
	// postgresql://ccvsuser:ccvspwd@localhost/ccvs
	db, err := postgres.New(os.Getenv("DATABASE_URL"), cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	v1 := e.Group("/container-scanning")

	at.NewHTTP(al.New(analysis.Initialize(db, cfg.Vendors), log), v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
