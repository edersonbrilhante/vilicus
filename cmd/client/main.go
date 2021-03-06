package main

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/edersonbrilhante/vilicus/pkg/client"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
)

func main() {

	cfgPath := flag.String("p", "./configs/conf.docker-compose.yaml", "Path to config file")
	imgs := flag.String("i", "", "Comma-separated list of images")
	template := flag.String("t", "", "Output template")
	output := flag.String("o", "", "Output File")
	flag.Parse()

	if *imgs == "" {
		checkErr(errors.New("No image was passed"))
	}

	imgList := strings.Split(*imgs, ",")

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	out := os.Stdout
	if *output != "" {
		if out, err = os.Create(*output); err != nil {
			panic(err.Error())
		}
	}

	checkErr(client.Start(cfg, imgList, *template, out))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
