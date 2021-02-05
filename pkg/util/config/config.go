package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Load returns Configuration struct
func Load(path string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}
	var cfg = new(Configuration)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}

// Configuration holds data necessary for configuring application
type Configuration struct {
	Server  *Server   `yaml:"server,omitempty"`
	DB      *Database `yaml:"database,omitempty"`
	Vendors *Vendors  `yaml:"vendors,omitempty"`
}

// Database holds data necessary for database configuration
type Database struct {
	LogQueries bool `yaml:"log_queries,omitempty"`
	Timeout    int  `yaml:"timeout_seconds,omitempty"`
}

// Server holds data necessary for server configuration
type Server struct {
	Port         string `yaml:"port,omitempty"`
	Debug        bool   `yaml:"debug,omitempty"`
	ReadTimeout  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int    `yaml:"write_timeout_seconds,omitempty"`
}

// Vendors is the struct that stores a list of vendors
type Vendors struct {
	Clair         *Clair         `yaml:"clair,omitempty"`
	AnchoreEngine *AnchoreEngine `yaml:"anchore_engine,omitempty"`
	Trivy         *Trivy         `yaml:"trivy,omitempty"`
}

// Clair is the struct that stores a vendor and its config
type Clair struct {
	URL string `yaml:"url"`
}

// AnchoreEngine is the struct that stores a vendor and its config
type AnchoreEngine struct {
	URL       string `yaml:"url"`
	User      string `yaml:"user"`
	Pass      string `yaml:"pass"`
	AsAccount string `yaml:"as_account"`
}

// Trivy is the struct that stores a vendor and its config
type Trivy struct {
	URL     string `yaml:"url"`
	Timeout int    `yaml:"timeout_seconds"`
}
