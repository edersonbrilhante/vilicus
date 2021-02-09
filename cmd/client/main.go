package main

import (
	"errors"
	"flag"
	"strings"

	"github.com/edersonbrilhante/vilicus/pkg/client"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
)

func main() {

	cfgPath := flag.String("p", "./configs/conf.docker-compose.yaml", "Path to config file")
	imgs := flag.String("i", "", "Comma-separated list of images")
	flag.Parse()

	if *imgs == "" {
		checkErr(errors.New("No image was passed"))
	}

	imgList := strings.Split(*imgs, ",")

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	checkErr(client.Start(cfg, imgList))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
