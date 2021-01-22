package ccvs

import (
	"fmt"
	"time"
)

// Analysis is the struct that stores all data from analysis performed.
type Analysis struct {
	ID        string    `json:"id,omitempty" pg:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Image     string    `json:"image,omitempty" pg:"image"`
	Status    string    `json:"status,omitempty" pg:"status"`
	CreatedAt time.Time `json:"created_at,omitempty" pg:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" pg:"updated_at"`
	Result    string    `json:"result,omitempty" pg:"result"`
	Errors    []string  `json:"errors,omitempty" pg:"errors"`
	Results   Results   `json:"ccvs_results,omitempty" pg:"ccvs_results"`
}

// Results is a struct that represents ccvs scan results.
type Results struct {
	ClairResult         VendorResults `json:"clair,omitempty"`
	AnchoreEngineResult VendorResults `json:"anchore_engine,omitempty"`
	TrivyResult         VendorResults `json:"trivy,omitempty"`
}

// VendorResults stores all Unknown, Negligible Low, Medium, High and Critical vulnerabilities for a vendor
type VendorResults struct {
	UnknownVulns    []Vuln `json:"unknown_vulns,omitempty"`
	NegligibleVulns []Vuln `json:"negligible_vulns,omitempty"`
	LowVulns        []Vuln `json:"low_vulns,omitempty"`
	MediumVulns     []Vuln `json:"medium_vulns,omitempty"`
	HighVulns       []Vuln `json:"high_vulns,omitempty"`
	CriticalVulns   []Vuln `json:"critical_vulns,omitempty"`
}

func (v VendorResults) String() string {
	return fmt.Sprintf(
		"VendorResults<Unknown:%d | Negligible:%d | Low:%d | Medium:%d | High:%d | Critical%d>",
		len(v.UnknownVulns), len(v.NegligibleVulns), len(v.LowVulns),
		len(v.MediumVulns), len(v.HighVulns), len(v.CriticalVulns),
	)
}

// Vuln is the struct that stores vulnerability information.
type Vuln struct {
	Fix            string `json:"fix"`
	URL            string `json:"url"`
	Name           string `json:"name"`
	Severity       string `json:"severity"`
	PackageName    string `json:"package_name"`
	PackageVersion string `json:"package_version"`
}

func (v Vuln) String() string {
	return fmt.Sprintf(
		"Vuln<[%s][%s] %s(%s) - Fix:%s>",
		v.Severity, v.Name, v.PackageName, v.PackageVersion, v.Fix,
	)
}
