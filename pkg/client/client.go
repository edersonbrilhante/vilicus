package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/edersonbrilhante/vilicus/pkg/types"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
	"github.com/go-resty/resty/v2"
)

type client struct {
	config *config.Client
}

func Start(cfg *config.Configuration, imgList []string) error {
	cli := client{config: cfg.Client}

	for _, img := range imgList {
		fmt.Printf("Starting analysis with image %s \n", img)
		analysis, err := cli.createAnalysis(img)
		if err != nil {
			return err
		}
		analysis, err = cli.monitorResult(analysis.ID)
		if err != nil {
			return err
		}
		fmt.Print(analysis.Results)
		fmt.Printf("Finised analysis with image %s \n", img)
		fmt.Printf("Results: %s \n", analysis.Results)

	}
	return nil
}

func (c client) createAnalysis(img string) (types.Analysis, error) {
	client := resty.New()
	analysis := types.Analysis{}

	resp, err := client.R().
		SetBody(types.Analysis{Image: img}).
		Post(c.config.URL + "/analysis")
	if err != nil {
		return analysis, err
	}
	err = json.Unmarshal(resp.Body(), &analysis)

	return analysis, err

}

func (c client) monitorResult(rid string) (types.Analysis, error) {
	analysis := types.Analysis{}
	timeout := time.After(60 * time.Minute)
	retryTick := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-timeout:
			return analysis, errors.New("time out")
		case <-retryTick.C:
			analysis, err := c.getAnalysis(rid)
			if err != nil {
				return analysis, err
			}
			if analysis.Status == "finished" {
				return analysis, nil
			}
		}
	}

}

func (c client) getAnalysis(rid string) (types.Analysis, error) {
	client := resty.New()
	analysis := types.Analysis{}

	resp, err := client.R().
		Get(c.config.URL + "/analysis/" + rid)
	if err != nil {
		return analysis, err
	}
	err = json.Unmarshal(resp.Body(), &analysis)

	return analysis, err

}
