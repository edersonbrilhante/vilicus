package vendor

import (
	"github.com/edersonbrilhante/ccvs"
	"github.com/edersonbrilhante/ccvs/pkg/util/config"
)

type ImageAnalysisRequest struct {
	source ImageSource `json:"source"`
}

type ImageSource struct {
	tag RegistryTagSource `json:"tag"`
}

type RegistryTagSource struct {
	pullstring string `json:"pullstring"`
	dockerfile string `json:"dockerfile"`
}

type Anchore struct {
	config    *config.AnchoreEngine
	resultRaw []string
	analysis  *ccvs.Analysis
}

func (a *Anchore) Analyzer(al *ccvs.Analysis) error {
	a.analysis = al
	return nil
}

func (a *Anchore) Parser() error {
	r := ccvs.VendorResults{}

	v := ccvs.Vuln{
		Fix:            "",
		URL:            "",
		Name:           "",
		Severity:       "unknown",
		PackageName:    "",
		PackageVersion: "",
	}
	switch v.Severity {
	case "unknown":
		r.UnknownVulns = append(r.UnknownVulns, v)
	case "negligible":
		r.NegligibleVulns = append(r.NegligibleVulns, v)
	case "low":
		r.LowVulns = append(r.LowVulns, v)
	case "medium":
		r.MediumVulns = append(r.MediumVulns, v)
	case "high":
		r.HighVulns = append(r.HighVulns, v)
	case "critical":
		r.CriticalVulns = append(r.CriticalVulns, v)
	}

	a.analysis.Results.AnchoreEngineResult = r

	return nil
}
