package clair

import (
	"errors"
	"fmt"
	"time"

	// "github.com/quay/claircore"

	"github.com/edersonbrilhante/ccvs"
	"github.com/edersonbrilhante/ccvs/pkg/util/config"
	"github.com/go-resty/resty/v2"

	"context"
	"encoding/json"
	"strings"

	"net/http"
	"sync"
	// "github.com/quay/clair/v4/httptransport/client"
	// "github.com/optiopay/klar/clair"
	// "github.com/optiopay/klar/docker"
)

type Clair struct {
	Config    *config.Clair
	resultRaw *VulnerabilityReport
	analysis  *ccvs.Analysis
}

var (
	rtMu  sync.Mutex
	rtMap = map[string]http.RoundTripper{}
)

func (a *Clair) Analyzer(al *ccvs.Analysis) error {
	a.analysis = al
	imgResp, err := a.addImage()
	if err != nil {
		return err
	}
	_, err = a.checkAnalysisStatus(imgResp)
	if err != nil {
		return err
	}

	vulnResp, err := a.getVuln(imgResp)
	if err != nil {
		return err
	}
	a.resultRaw = &vulnResp

	return nil
}

// image, err := docker.NewImage(&docker.DockerConfig{
// 	ImageName: al.Image
// })
// if err != nil {
// 	return errors.New("Can't parse qname: %s", err)
// }

// err = image.Pull()
// if err != nil {
// 	return errors.New("Can't pull image: %s", err)
// }

// var vs []*clair.Vulnerability
// c := clair.NewClair(conf.ClairAddr, 3, conf.ClairTimeout)
// vs, err = c.Analyse(image)
// if err != nil {
// 	return errors.New("Failed to analyze using API: %s", err)
// }
// a.resultRaw = vs

func (a *Clair) addImage() (IndexReport, error) {

	imgResp := IndexReport{}

	ctx := context.Background()
	manifest, err := Inspect(ctx, a.analysis.Image)
	if err != nil {
		return imgResp, err
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(manifest).
		Post(a.Config.URL + "/indexer/api/v1/index_report")
	if err != nil {
		return imgResp, err
	}
	err = json.Unmarshal(resp.Body(), &imgResp)

	return imgResp, err
}

func (a *Clair) checkAnalysisStatus(imgResp IndexReport) (bool, error) {
	if imgResp.State == "IndexFinished" {
		return true, nil
	}

	fmt.Println("Sleeping 10s")
	time.Sleep(10000 * time.Millisecond) //TODO:use select instead sleep

	client := resty.New()

	newImgResp := IndexReport{}
	resp, err := client.R().
		Get(a.Config.URL + "/indexer/api/v1/index_report/" + imgResp.ManifestHash)

	if err == nil {
		err = json.Unmarshal(resp.Body(), &newImgResp)
	}

	if err != nil {
		return false, err
	}

	return a.checkAnalysisStatus(newImgResp)
}

func (a *Clair) getVuln(imgResp IndexReport) (VulnerabilityReport, error) {
	client := resty.New()

	vulnResp := VulnerabilityReport{}
	resp, err := client.R().
		Get(a.Config.URL + "/matcher/api/v1/vulnerability_report/" + imgResp.ManifestHash)
	if err != nil {
		return vulnResp, err
	}

	err = json.Unmarshal(resp.Body(), &vulnResp)

	return vulnResp, err
}

func (a *Clair) Parser() error {
	if a.resultRaw == nil {
		return errors.New("Result is empty")
	}
	fmt.Printf("%v: found %v vulns\n", a.analysis.Image, len(a.resultRaw.Vulnerabilities))
	r := ccvs.VendorResults{}

	for p, vis := range a.resultRaw.PackageVulnerabilities {
		for _, vi := range vis {

			pkg := a.resultRaw.Packages[p]

			v := a.resultRaw.Vulnerabilities[vi]

			vuln := ccvs.Vuln{
				Fix:            v.FixedInVersion,
				URL:            strings.ReplaceAll(v.Links, " ", ", "),
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

	a.analysis.Results.ClairResult = r

	return nil
}
