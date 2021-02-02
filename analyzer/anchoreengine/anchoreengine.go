package anchoreengine

import (
	"errors"
	"fmt"
	"time"

	"github.com/edersonbrilhante/vilicus"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
	"github.com/go-resty/resty/v2"
)

type ImageAnalysisRequest struct {
	Source ImageSource `json:"source"`
}

type ImageSource struct {
	Tag RegistryTagSource `json:"tag"`
}
type RegistryTagSource struct {
	Pullstring string `json:"pullstring"`
}

type AnchoreImage struct {
	AnalysisStatus string        `json:"analysis_status"`
	AnalyzedAt     interface{}   `json:"analyzed_at"`
	Annotations    Annotations   `json:"annotations"`
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
type Annotations struct {
}
type Metadata struct {
	Arch           interface{} `json:"arch"`
	Distro         interface{} `json:"distro"`
	DistroVersion  interface{} `json:"distro_version"`
	DockerfileMode interface{} `json:"dockerfile_mode"`
	ImageSize      interface{} `json:"image_size"`
	LayerCount     interface{} `json:"layer_count"`
}
type ImageContent struct {
	Metadata Metadata `json:"metadata"`
}
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

type VulnerabilityResponse struct {
	ImageDigest       string            `json:"imageDigest"`
	Vulnerabilities   []Vulnerabilities `json:"vulnerabilities"`
	VulnerabilityType string            `json:"vulnerability_type"`
}
type CvssV2 struct {
	BaseScore           float64 `json:"base_score"`
	ExploitabilityScore float64 `json:"exploitability_score"`
	ImpactScore         float64 `json:"impact_score"`
}
type CvssV3 struct {
	BaseScore           float64 `json:"base_score"`
	ExploitabilityScore float64 `json:"exploitability_score"`
	ImpactScore         float64 `json:"impact_score"`
}
type NvdData struct {
	CvssV2 CvssV2 `json:"cvss_v2"`
	CvssV3 CvssV3 `json:"cvss_v3"`
	ID     string `json:"id"`
}
type Vulnerabilities struct {
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

type Anchore struct {
	Config    *config.AnchoreEngine
	resultRaw *VulnerabilityResponse
	analysis  *vilicus.Analysis
}

func (a *Anchore) Analyzer(al *vilicus.Analysis) error {
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

func (a *Anchore) Parser() error {
	if a.resultRaw == nil {
		return errors.New("Result is empty")
	}
	r := vilicus.VendorResults{}
	for _, v := range a.resultRaw.Vulnerabilities {

		vuln := vilicus.Vuln{
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
