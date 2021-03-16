package anchoreengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/edersonbrilhante/vilicus/pkg/types"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
	"github.com/go-resty/resty/v2"
)

// Anchore stores pointers to config and results from scan.
type Anchore struct {
	Config    *config.AnchoreEngine
	resultRaw *vulnerabilityResponse
	analysis  *types.Analysis
}

// Analyzer runs an analysis and stores the results in Anchore.resultRaw.
func (a *Anchore) Analyzer(al *types.Analysis) error {
	a.analysis = al
	imgResp, err := a.addImage()
	if err != nil {
		return err
	}

	_, err = a.monitorResult(imgResp[0].ImageDigest)
	if err != nil {
		return err
	}

	vulnResp, err := a.getVuln(imgResp[0].ImageDigest)
	if err != nil {
		return err
	}
	a.resultRaw = &vulnResp

	return nil
}

// Parser parses Anchore.resultRaw and store the final data into a type Analysis.
func (a *Anchore) Parser() error {
	if a.resultRaw == nil {
		return errors.New("Result is empty")
	}
	r := types.VendorResults{}
	for _, v := range a.resultRaw.Vulnerabilities {

		vuln := types.Vuln{
			Fix:            v.Fix,
			URL:            filterValidURLs([]string{v.URL}),
			Name:           v.Vuln,
			Vendor:         "AnchoreEngine",
			Severity:       v.Severity,
			PackageName:    v.PackageName,
			PackageVersion: v.PackageVersion,
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

	a.analysis.Results.AnchoreEngineResult = r

	return nil
}

func (a *Anchore) addImage() ([]anchoreImage, error) {
	client := resty.New()
	client.SetBasicAuth(a.Config.User, a.Config.Pass)

	imgResp := []anchoreImage{}

	resp, err := client.R().
		SetBody(imageAnalysisRequest{Source: imageSource{Tag: registryTagSource{Pullstring: a.analysis.Image}}}).
		SetResult(&imgResp).
		Post(a.Config.URL + "/images?force=true&autosubscribe=false")
	if err != nil {
		return imgResp, err
	}

	if len(imgResp) == 0 {
		return imgResp, errors.New(resp.String())
	}
	return imgResp, err
}

func (a *Anchore) monitorResult(rid string) (anchoreImage, error) {
	aimage := anchoreImage{}
	timeout := time.After(60 * time.Minute)
	retryTick := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-timeout:
			return aimage, errors.New("time out")
		case <-retryTick.C:
			fmt.Println(aimage.AnalysisStatus)
			aimage, err := a.getAnalysis(rid)
			if err != nil {
				return aimage, err
			}
			if aimage.AnalysisStatus == "analyzed" {
				return aimage, nil
			}
			if aimage.AnalysisStatus == "analysis_failed" {
				return aimage, errors.New("analysis failed")
			}
		}
	}

}

func (a *Anchore) getAnalysis(imageDigest string) (anchoreImage, error) {
	client := resty.New()
	client.SetBasicAuth(a.Config.User, a.Config.Pass)
	aimage := []anchoreImage{}

	resp, err := client.R().
		Get(a.Config.URL + "/images/" + imageDigest)
	if err != nil {
		return aimage[0], err
	}
	err = json.Unmarshal(resp.Body(), &aimage)

	return aimage[0], err

}

func (a *Anchore) getVuln(imageDigest string) (vulnerabilityResponse, error) {
	client := resty.New()
	client.SetBasicAuth(a.Config.User, a.Config.Pass)
	vuln := vulnerabilityResponse{}

	resp, err := client.R().
		Get(a.Config.URL + "/images/" + imageDigest + "/vuln/all")

	if err != nil {
		return vuln, err
	}
	err = json.Unmarshal(resp.Body(), &vuln)

	return vuln, err
}

func filterValidURLs(urls []string) []string {
	validURLs := []string{}
	for _, ur := range urls {
		u, err := url.Parse(ur)
		if err == nil && u.Scheme != "" && u.Host != "" {
			validURLs = append(validURLs, u.String())
		}
	}
	return validURLs
}
