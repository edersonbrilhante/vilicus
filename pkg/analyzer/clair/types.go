package clair

import (
	"time"
)

// IndexReport is an intermediate data structure describing the contents of a container image
type IndexReport struct {
	ManifestHash  string                   `json:"manifest_hash"`
	State         string                   `json:"state"`
	Packages      map[string]Package       `json:"packages"`
	Distributions map[string]Distribution  `json:"distributions"`
	Environments  map[string][]Environment `json:"environments"`
	Success       bool                     `json:"success"`
	Err           string                   `json:"err"`
}

// SourcePackage stores a source package affiliated with a Package
type SourcePackage struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	Kind              string `json:"kind"`
	NormalizedVersion string `json:"normalized_version"`
	Cpe               string `json:"cpe"`
}

// Package storess an item discovered by indexing a Manifest"
type Package struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Version           string        `json:"version"`
	Kind              string        `json:"kind"`
	SourcePackage     SourcePackage `json:"source"`
	NormalizedVersion string        `json:"normalized_version"`
	Arch              string        `json:"arch"`
	Cpe               string        `json:"cpe"`
}

// Distribution stores an indexed distribution discovered in a layer.
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

// Environment stores a particular package was discovered in.
type Environment struct {
	PackageDb      string `json:"package_db"`
	IntroducedIn   string `json:"introduced_in"`
	DistributionID string `json:"distribution_id"`
}

// VulnerabilityReport stores discovered packages, package environments, and package vulnerabilities within a Manifest.
type VulnerabilityReport struct {
	ManifestHash           string                   `json:"manifest_hash"`
	Packages               map[string]Package       `json:"packages"`
	Repository             Repository               `json:"repository"`
	Environments           map[string][]Environment `json:"environments"`
	Vulnerabilities        map[string]Vulnerability `json:"vulnerabilities"`
	PackageVulnerabilities map[string][]string      `json:"package_vulnerabilities"`
}

// Repository stores a package repository
type Repository struct {
	Cpe  string `json:"cpe"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	URI  string `json:"uri"`
}

// Vulnerability an unique item indexed by Clair
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
