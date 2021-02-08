package anchoreengine

import (
	"errors"
	"fmt"
	"time"

	"github.com/edersonbrilhante/vilicus/pkg/types"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
	"github.com/go-resty/resty/v2"
)

// ImageAnalysisRequest stores request to add an image to be watched and analyzed by the engine.
type ImageAnalysisRequest struct {
	Source ImageSource `json:"source"`
}

// ImageSource stores set of analysis source types.
type ImageSource struct {
	Tag RegistryTagSource `json:"tag"`
}

// RegistryTagSource stores an image reference using a tag in a registry, this is the most common source type.
type RegistryTagSource struct {
	Pullstring string `json:"pullstring"`
}

// AnchoreImage stores information about an image analysis.
type AnchoreImage struct {
	AnalysisStatus string        `json:"analysis_status"`
	AnalyzedAt     interface{}   `json:"analyzed_at"`
	Annotations    interface{}   `json:"annotations"`
	CreatedAt      time.Time     `json:"created_at"`
	ImageDigest    string        `json:"imageDigest"`
	ImageContent   ImageContent  `json:"image_content"`
	ImageDetail    []ImageDetail `json:"image_detail"`
	ImageStatus    string        `json:"image_status"`
	ImageType      string        `json:"image_type"`
	LastUpdated    time.Time     `json:"last_updated"`
	ParentDigest   string        `json:"parentDigest"`
	UserID         string        `json:"userId"`
}

// Metadata stores content record for a specific image, containing different content type entries.
type Metadata struct {
	Arch           interface{} `json:"arch"`
	Distro         interface{} `json:"distro"`
	DistroVersion  interface{} `json:"distro_version"`
	DockerfileMode interface{} `json:"dockerfile_mode"`
	ImageSize      interface{} `json:"image_size"`
	LayerCount     interface{} `json:"layer_count"`
}

// ImageContent stores a metadata struct
type ImageContent struct {
	Metadata Metadata `json:"metadata"`
}

// ImageDetail stores a metadata detail record for a specific image.
type ImageDetail struct {
	CreatedAt     time.Time   `json:"created_at"`
	Digest        string      `json:"digest"`
	Dockerfile    interface{} `json:"dockerfile"`
	Fulldigest    string      `json:"fulldigest"`
	Fulltag       string      `json:"fulltag"`
	ImageDigest   string      `json:"imageDigest"`
	ImageID       string      `json:"imageId"`
	LastUpdated   time.Time   `json:"last_updated"`
	Registry      string      `json:"registry"`
	Repo          string      `json:"repo"`
	Tag           string      `json:"tag"`
	TagDetectedAt time.Time   `json:"tag_detected_at"`
	UserID        string      `json:"userId"`
}

// VulnerabilityResponse envelope containing list of vulnerabilities.
type VulnerabilityResponse struct {
	ImageDigest       string          `json:"imageDigest"`
	Vulnerabilities   []Vulnerability `json:"vulnerabilities"`
	VulnerabilityType string          `json:"vulnerability_type"`
}

// CVSSV2 stores information about vulnerability with CVSSV2 data.
type CVSSV2 struct {
	BaseScore           float64 `json:"base_score"`
	ExploitabilityScore float64 `json:"exploitability_score"`
	ImpactScore         float64 `json:"impact_score"`
}

// CVSSV3 stores information about vulnerability with CVSSV3 data.
type CVSSV3 struct {
	BaseScore           float64 `json:"base_score"`
	ExploitabilityScore float64 `json:"exploitability_score"`
	ImpactScore         float64 `json:"impact_score"`
}

// NvdData stores information about Nvd Data item.
type NvdData struct {
	CVSSV2 CVSSV2 `json:"cvss_v2"`
	CVSSV3 CVSSV3 `json:"cvss_v3"`
	ID     string `json:"id"`
}

// Vulnerability stores a information about vulnerability.
type Vulnerability struct {
	Feed           string        `json:"feed"`
	FeedGroup      string        `json:"feed_group"`
	Fix            string        `json:"fix"`
	NvdData        []NvdData     `json:"nvd_data"`
	Package        string        `json:"package"`
	PackageCpe     string        `json:"package_cpe"`
	PackageCpe23   string        `json:"package_cpe23"`
	PackageName    string        `json:"package_name"`
	PackagePath    string        `json:"package_path"`
	PackageType    string        `json:"package_type"`
	PackageVersion string        `json:"package_version"`
	Severity       string        `json:"severity"`
	URL            string        `json:"url"`
	VendorData     []interface{} `json:"vendor_data"`
	Vuln           string        `json:"vuln"`
}

// Anchore stores pointers to config and results from scan.
type Anchore struct {
	Config    *config.AnchoreEngine
	resultRaw *VulnerabilityResponse
	analysis  *types.Analysis
}

// Analyzer runs an analysis and stores the results in Anchore.resultRaw.
func (a *Anchore) Analyzer(al *types.Analysis) error {
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

// Parser parses Anchore.resultRaw and store the final data into a type Analysis.
func (a *Anchore) Parser() error {
	if a.resultRaw == nil {
		return errors.New("Result is empty")
	}
	r := types.VendorResults{}
	for _, v := range a.resultRaw.Vulnerabilities {

		vuln := types.Vuln{
			Fix:            v.Fix,
			URL:            []string{v.URL},
			Name:           v.Vuln,
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

func (a *Anchore) addImage() ([]AnchoreImage, error) {
	client := resty.New()
	client.SetBasicAuth(a.Config.User, a.Config.Pass)

	imgResp := []AnchoreImage{}

	resp, err := client.R().
		SetBody(ImageAnalysisRequest{Source: ImageSource{Tag: RegistryTagSource{Pullstring: a.analysis.Image}}}).
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

func (a *Anchore) checkAnalysisStatus(imgResp []AnchoreImage) (bool, error) {
	if imgResp[0].AnalysisStatus == "analyzed" {
		return true, nil
	}

	fmt.Println("Sleeping 10s")
	time.Sleep(10000 * time.Millisecond) //TODO:use select instead sleep

	client := resty.New()
	client.SetBasicAuth(a.Config.User, a.Config.Pass)

	newImgResp := []AnchoreImage{}
	_, err := client.R().
		SetResult(&newImgResp).
		Get(a.Config.URL + "/images/" + imgResp[0].ImageDigest)

	if err != nil {
		return false, err
	}

	return a.checkAnalysisStatus(newImgResp)
}

func (a *Anchore) getVuln(imgResp []AnchoreImage) (VulnerabilityResponse, error) {
	client := resty.New()
	client.SetBasicAuth(a.Config.User, a.Config.Pass)

	vulnResp := VulnerabilityResponse{}
	_, err := client.R().
		SetResult(&vulnResp).
		Get(a.Config.URL + "/images/" + imgResp[0].ImageDigest + "/vuln/all")

	return vulnResp, err
}
