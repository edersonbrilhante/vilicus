package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/edersonbrilhante/vilicus/pkg/report"
	"github.com/edersonbrilhante/vilicus/pkg/types"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
	"github.com/go-resty/resty/v2"
	"golang.org/x/xerrors"
)

type client struct {
	config *config.Client
}

// Start starts the scanning image
func Start(cfg *config.Configuration, imgList []string, template string, out io.Writer) error {
	var wg sync.WaitGroup
	resCh := make(chan types.Analysis, len(imgList))
	errCh := make(chan error, len(imgList))

	cli := client{config: cfg.Client}
	cli.processImages(imgList, resCh, errCh, &wg)
	err := cli.reportResult(resCh, errCh, template, out)

	return err
}

func (c client) reportResult(resCh chan types.Analysis, errCh chan error, template string, out io.Writer) error {

	var analyses types.Analyses

	for r := range resCh {
		analyses = append(analyses, r)
	}

	for e := range errCh {
		fmt.Printf("Error in analysis: %s\n", e)
	}

	err := report.WriteAnalyses(out, analyses, template)

	return err
}

func (c client) processImages(imgList []string, resCh chan types.Analysis, errCh chan error, wg *sync.WaitGroup) {
	for _, img := range imgList {
		wg.Add(1)
		go c.scan(img, resCh, errCh, wg)
	}
	wg.Wait()
	close(errCh)
	close(resCh)
}

func (c client) scan(img string, resCh chan types.Analysis, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	analysis, err := c.createAnalysis(img)
	if err != nil {
		errCh <- xerrors.Errorf("error creating analysis with image %s: %w \n", img, err)
		return
	}

	analysis, err = c.monitorResult(analysis.ID)
	if err != nil {
		errCh <- xerrors.Errorf("error getting result from analysis with image %s: %s \n", img, err)
		return
	}

	resCh <- analysis
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
