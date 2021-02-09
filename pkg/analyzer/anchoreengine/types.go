package anchoreengine

import (
	"time"
)

type imageAnalysisRequest struct {
	Source imageSource `json:"source"`
}

type imageSource struct {
	Tag registryTagSource `json:"tag"`
}

type registryTagSource struct {
	Pullstring string `json:"pullstring"`
}

type anchoreImage struct {
	AnalysisStatus string        `json:"analysis_status"`
	AnalyzedAt     interface{}   `json:"analyzed_at"`
	Annotations    interface{}   `json:"annotations"`
	CreatedAt      time.Time     `json:"created_at"`
	ImageDigest    string        `json:"imageDigest"`
	ImageContent   imageContent  `json:"image_content"`
	ImageDetail    []imageDetail `json:"image_detail"`
	ImageStatus    string        `json:"image_status"`
	ImageType      string        `json:"image_type"`
	LastUpdated    time.Time     `json:"last_updated"`
	ParentDigest   string        `json:"parentDigest"`
	UserID         string        `json:"userId"`
}

type metadata struct {
	Arch           interface{} `json:"arch"`
	Distro         interface{} `json:"distro"`
	DistroVersion  interface{} `json:"distro_version"`
	DockerfileMode interface{} `json:"dockerfile_mode"`
	ImageSize      interface{} `json:"image_size"`
	LayerCount     interface{} `json:"layer_count"`
}

type imageContent struct {
	Metadata metadata `json:"metadata"`
}

type imageDetail struct {
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

type vulnerabilityResponse struct {
	ImageDigest       string          `json:"imageDigest"`
	Vulnerabilities   []vulnerability `json:"vulnerabilities"`
	VulnerabilityType string          `json:"vulnerability_type"`
}

type cvssV2 struct {
	BaseScore           float64 `json:"base_score"`
	ExploitabilityScore float64 `json:"exploitability_score"`
	ImpactScore         float64 `json:"impact_score"`
}

type cvssV3 struct {
	BaseScore           float64 `json:"base_score"`
	ExploitabilityScore float64 `json:"exploitability_score"`
	ImpactScore         float64 `json:"impact_score"`
}

type nvdData struct {
	CVSSV2 cvssV2 `json:"cvss_v2"`
	CVSSV3 cvssV3 `json:"cvss_v3"`
	ID     string `json:"id"`
}

type vulnerability struct {
	Feed           string        `json:"feed"`
	FeedGroup      string        `json:"feed_group"`
	Fix            string        `json:"fix"`
	NvdData        []nvdData     `json:"nvd_data"`
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
