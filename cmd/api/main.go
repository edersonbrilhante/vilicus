package main

import (
	"flag"

	"github.com/edersonbrilhante/vilicus/pkg/api"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
)

func main() {

	cfgPath := flag.String("p", "./configs/conf.docker-compose.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
