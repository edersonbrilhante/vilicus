package clair

import (
	"time"
)

type indexReport struct {
	ManifestHash  string                   `json:"manifest_hash"`
	State         string                   `json:"state"`
	Packages      map[string]pkg           `json:"packages"`
	Distributions map[string]distribution  `json:"distributions"`
	Environments  map[string][]environment `json:"environments"`
	Success       bool                     `json:"success"`
	Err           string                   `json:"err"`
}

type sourcePackage struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	Kind              string `json:"kind"`
	NormalizedVersion string `json:"normalized_version"`
	Cpe               string `json:"cpe"`
}

type pkg struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Version           string        `json:"version"`
	Kind              string        `json:"kind"`
	SourcePackage     sourcePackage `json:"source"`
	NormalizedVersion string        `json:"normalized_version"`
	Arch              string        `json:"arch"`
	Cpe               string        `json:"cpe"`
}

type distribution struct {
	ID              string `json:"id"`
	Did             string `json:"did"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	VersionCodeName string `json:"version_code_name"`
	VersionID       string `json:"version_id"`
	Arch            string `json:"arch"`
	Cpe             string `json:"cpe"`
	PrettyName      string `json:"pretty_name"`
}

type environment struct {
	PackageDb      string `json:"package_db"`
	IntroducedIn   string `json:"introduced_in"`
	DistributionID string `json:"distribution_id"`
}

type vulnerabilityReport struct {
	ManifestHash           string                   `json:"manifest_hash"`
	Packages               map[string]pkg           `json:"packages"`
	Repository             repository               `json:"repository"`
	Environments           map[string][]environment `json:"environments"`
	Vulnerabilities        map[string]vulnerability `json:"vulnerabilities"`
	PackageVulnerabilities map[string][]string      `json:"package_vulnerabilities"`
}

type repository struct {
	Cpe  string `json:"cpe"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	URI  string `json:"uri"`
}

type vulnerability struct {
	ID                 string       `json:"id"`
	Updater            string       `json:"updater"`
	Name               string       `json:"name"`
	Description        string       `json:"description"`
	Issued             time.Time    `json:"issued"`
	Links              string       `json:"links"`
	Severity           string       `json:"severity"`
	NormalizedSeverity string       `json:"normalized_severity"`
	Package            pkg          `json:"package"`
	Distribution       distribution `json:"distribution"`
	Repository         repository   `json:"repository"`
	FixedInVersion     string       `json:"fixed_in_version"`
}
