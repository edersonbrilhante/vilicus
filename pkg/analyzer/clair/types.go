package clair

import (
	"time"
)

type IndexReport struct {
	ManifestHash  string                   `json:"manifest_hash"`
	State         string                   `json:"state"`
	Packages      map[string]Package       `json:"packages"`
	Distributions map[string]Distribution  `json:"distributions"`
	Environments  map[string][]Environment `json:"environments"`
	Success       bool                     `json:"success"`
	Err           string                   `json:"err"`
}
type Source struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	Kind              string `json:"kind"`
	NormalizedVersion string `json:"normalized_version"`
	Cpe               string `json:"cpe"`
}
type Package struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	Kind              string `json:"kind"`
	Source            Source `json:"source"`
	NormalizedVersion string `json:"normalized_version"`
	Arch              string `json:"arch"`
	Cpe               string `json:"cpe"`
}

type Distribution struct {
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
type Environment struct {
	PackageDb      string `json:"package_db"`
	IntroducedIn   string `json:"introduced_in"`
	DistributionID string `json:"distribution_id"`
	// RepositoryIds  interface{} `json:"repository_ids"`
}

type VulnerabilityReport struct {
	ManifestHash string             `json:"manifest_hash"`
	Packages     map[string]Package `json:"packages"`
	// Distributions          map[string]Distribution  `json:"distributions"`
	Repository             Repository               `json:"repository"`
	Environments           map[string][]Environment `json:"environments"`
	Vulnerabilities        map[string]Vulnerability `json:"vulnerabilities"`
	PackageVulnerabilities map[string][]string      `json:"package_vulnerabilities"`
}

type Repository struct {
	Cpe  string `json:"cpe"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	Uri  string `json:"uri"`
}
type Vulnerability struct {
	ID                 string       `json:"id"`
	Updater            string       `json:"updater"`
	Name               string       `json:"name"`
	Description        string       `json:"description"`
	Issued             time.Time    `json:"issued"`
	Links              string       `json:"links"`
	Severity           string       `json:"severity"`
	NormalizedSeverity string       `json:"normalized_severity"`
	Package            Package      `json:"package"`
	Distribution       Distribution `json:"distribution"`
	Repository         Repository   `json:"repository"`
	FixedInVersion     string       `json:"fixed_in_version"`
}
