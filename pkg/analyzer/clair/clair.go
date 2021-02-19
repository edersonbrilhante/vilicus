package clair

import (
	"errors"
	"fmt"
	"time"

	"github.com/edersonbrilhante/vilicus/pkg/types"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
	"github.com/go-resty/resty/v2"

	"context"
	"encoding/json"
	"strings"

	"net/http"
	"sync"
)

// Clair stores pointers to config and results from scan
type Clair struct {
	Config    *config.Clair
	resultRaw *vulnerabilityReport
	analysis  *types.Analysis
}

var (
	rtMu  sync.Mutex
	rtMap = map[string]http.RoundTripper{}
)

// Analyzer runs an analysis and stores the results in Clair.resultRaw
func (c *Clair) Analyzer(al *types.Analysis) error {
	c.analysis = al
	imgResp, err := c.addImage()
	if err != nil {
		return err
	}
	_, err = c.monitorResult(imgResp.ManifestHash)
	if err != nil {
		return err
	}

	vulnResp, err := c.getVuln(imgResp.ManifestHash)
	if err != nil {
		return err
	}
	c.resultRaw = &vulnResp

	return nil
}

// Parser parses Clair.resultRaw and store the final data into a type Analysis
func (c *Clair) Parser() error {
	if c.resultRaw == nil {
		return errors.New("Result is empty")
	}
	fmt.Printf("%v: found %v vulns\n", c.analysis.Image, len(c.resultRaw.Vulnerabilities))
	r := types.VendorResults{}

	for p, vis := range c.resultRaw.PackageVulnerabilities {
		for _, vi := range vis {

			pkg := c.resultRaw.Packages[p]

			v := c.resultRaw.Vulnerabilities[vi]

			vuln := types.Vuln{
				Fix:            v.FixedInVersion,
				URL:            strings.Split(v.Links, " "),
				Name:           v.Name,
				Severity:       v.NormalizedSeverity,
				PackageName:    pkg.Name,
				PackageVersion: pkg.Version,
			}
			switch vuln.Severity {
			case "Unknown":
				r.UnknownVulns = append(r.UnknownVulns, vuln)
			case "Negligible":
				r.NegligibleVulns = append(r.NegligibleVulns, vuln)
			case "Low":
				r.LowVulns = append(r.LowVulns, vuln)
			case "Medium":
				r.MediumVulns = append(r.MediumVulns, vuln)
			case "High":
				r.HighVulns = append(r.HighVulns, vuln)
			case "Critical":
				r.CriticalVulns = append(r.CriticalVulns, vuln)
			}
		}
	}

	c.analysis.Results.ClairResult = r

	return nil
}

func (c *Clair) addImage() (indexReport, error) {

	imgResp := indexReport{}

	ctx := context.Background()
	manifest, err := inspect(ctx, c.analysis.Image)
	if err != nil {
		return imgResp, err
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(manifest).
		Post(c.Config.URL + "/indexer/api/v1/index_report")
	if err != nil {
		return imgResp, err
	}
	err = json.Unmarshal(resp.Body(), &imgResp)

	return imgResp, err
}

func (c *Clair) monitorResult(rid string) (indexReport, error) {
	ireport := indexReport{}
	timeout := time.After(60 * time.Minute)
	retryTick := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-timeout:
			return ireport, errors.New("time out")
		case <-retryTick.C:
			fmt.Println(ireport.State)
			ireport, err := c.getAnalysis(rid)
			if err != nil {
				return ireport, err
			}
			if ireport.State == "IndexFinished" {
				return ireport, nil
			}
			if ireport.State == "IndexError" {
				return ireport, errors.New("analysis failed: " + ireport.Err)
			}
		}
	}

}

func (c *Clair) getAnalysis(rid string) (indexReport, error) {
	client := resty.New()
	ireport := indexReport{}

	resp, err := client.R().
		Get(c.Config.URL + "/indexer/api/v1/index_report/" + rid)
	if err != nil {
		return ireport, err
	}
	err = json.Unmarshal(resp.Body(), &ireport)

	return ireport, err

}

func (c *Clair) getVuln(rid string) (vulnerabilityReport, error) {
	client := resty.New()
	vulnResp := vulnerabilityReport{}

	resp, err := client.R().
		Get(c.Config.URL + "/matcher/api/v1/vulnerability_report/" + rid)
	if err != nil {
		return vulnResp, err
	}
	err = json.Unmarshal(resp.Body(), &vulnResp)

	return vulnResp, err
}
