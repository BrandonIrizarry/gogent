package yamlconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type YAMLConfig struct {
	MaxIterations int    `yaml:"max_iterations"`
	MaxFilesize   int    `yaml:"max_filesize"`
	WorkingDir    string `yaml:"working_dir"`
}

const (
	minMaxIterations = 20
	minMaxFilesize   = 1000

	// The working directory defaults to gogent's own current
	// working directory.
	defaultWorkingDir = "."
)

// NewYAMLConfig reads filename and returns a populated YAMLConfig
// struct plus an error. Validation is performed on fields where
// applicable.
func NewYAMLConfig(filename string) (YAMLConfig, error) {
	var cfg YAMLConfig

	yamlBytes, err := os.ReadFile(filename)
	if err != nil {
		return YAMLConfig{}, err
	}

	if err := yaml.Unmarshal(yamlBytes, &cfg); err != nil {
		return YAMLConfig{}, err
	}

	// Use default values in case values are missing or invalid.
	if cfg.MaxIterations < minMaxIterations {
		cfg.MaxIterations = minMaxIterations
	}

	if cfg.MaxFilesize < minMaxFilesize {
		cfg.MaxFilesize = minMaxFilesize
	}

	if cfg.WorkingDir == "" {
		cfg.WorkingDir = defaultWorkingDir
	}

	// Expand '~' prefix to the user's home directory.
	if cfg.WorkingDir[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return YAMLConfig{}, err
		}

		cfg.WorkingDir = filepath.Join(homeDir, cfg.WorkingDir[1:])
	}

	// Use the absolute path form of the working directory.
	cfg.WorkingDir, err = filepath.Abs(cfg.WorkingDir)
	if err != nil {
		return YAMLConfig{}, err
	}

	// Verify that cfg.WorkingDir is a directory.
	finfo, err := os.Stat(cfg.WorkingDir)
	if err != nil {
		return YAMLConfig{}, err
	}

	if ok := finfo.IsDir(); !ok {
		return YAMLConfig{}, fmt.Errorf("invalid working directory: %s", cfg.WorkingDir)
	}

	return cfg, nil
}
