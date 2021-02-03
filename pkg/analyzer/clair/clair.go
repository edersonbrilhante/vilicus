package clair

import (
	"errors"
	"fmt"
	"time"

	"github.com/edersonbrilhante/vilicus"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
	"github.com/go-resty/resty/v2"

	"context"
	"encoding/json"
	"strings"

	"net/http"
	"sync"
)

type Clair struct {
	Config    *config.Clair
	resultRaw *VulnerabilityReport
	analysis  *vilicus.Analysis
}

var (
	rtMu  sync.Mutex
	rtMap = map[string]http.RoundTripper{}
)

func (c *Clair) Analyzer(al *vilicus.Analysis) error {
	c.analysis = al
	imgResp, err := c.addImage()
	if err != nil {
		return err
	}
	_, err = c.checkAnalysisStatus(imgResp)
	if err != nil {
		return err
	}

	vulnResp, err := c.getVuln(imgResp)
	if err != nil {
		return err
	}
	c.resultRaw = &vulnResp

	return nil
}

func (c *Clair) addImage() (IndexReport, error) {

	imgResp := IndexReport{}

	ctx := context.Background()
	manifest, err := Inspect(ctx, c.analysis.Image)
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

func (c *Clair) checkAnalysisStatus(imgResp IndexReport) (bool, error) {
	if imgResp.State == "IndexFinished" {
		return true, nil
	}

	fmt.Println("Sleeping 10s")
	time.Sleep(10000 * time.Millisecond) //TODO:use select instead sleep

	client := resty.New()

	newImgResp := IndexReport{}
	resp, err := client.R().
		Get(c.Config.URL + "/indexer/api/v1/index_report/" + imgResp.ManifestHash)

	if err == nil {
		err = json.Unmarshal(resp.Body(), &newImgResp)
	}

	if err != nil {
		return false, err
	}

	return c.checkAnalysisStatus(newImgResp)
}

func (c *Clair) getVuln(imgResp IndexReport) (VulnerabilityReport, error) {
	client := resty.New()

	vulnResp := VulnerabilityReport{}
	resp, err := client.R().
		Get(c.Config.URL + "/matcher/api/v1/vulnerability_report/" + imgResp.ManifestHash)
	if err != nil {
		return vulnResp, err
	}

	err = json.Unmarshal(resp.Body(), &vulnResp)

	return vulnResp, err
}

func (c *Clair) Parser() error {
	if c.resultRaw == nil {
		return errors.New("Result is empty")
	}
	fmt.Printf("%v: found %v vulns\n", c.analysis.Image, len(c.resultRaw.Vulnerabilities))
	r := vilicus.VendorResults{}

	for p, vis := range c.resultRaw.PackageVulnerabilities {
		for _, vi := range vis {

			pkg := c.resultRaw.Packages[p]

			v := c.resultRaw.Vulnerabilities[vi]

			vuln := vilicus.Vuln{
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
