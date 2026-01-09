package yamlconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type YAMLConfig struct {
	MaxIterations int    `yaml:"max_iterations" validate:"required,gte=20,lte=100"`
	MaxFilesize   int    `yaml:"max_filesize" validate:"gte=1000,lte=200000"`
	WorkingDir    string `yaml:"working_dir" validate:"required"`
	RenderStyle   string `yaml:"render_style" validate:"oneof=light dark none"`
}

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

	// Validate the struct.
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("yaml"), ",", 2)[0]

		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}

		return name
	})

	if err := v.Struct(cfg); err != nil {
		var fieldErr validator.ValidationErrors
		var bld strings.Builder

		fmt.Fprintf(&bld, "Errors in %s:\n\n", filename)

		if errors.As(err, &fieldErr) {
			for _, verr := range fieldErr {
				value := verr.Value()
				field := verr.Field()

				fmt.Fprintf(&bld, "%s: invalid value: %v \n", field, value)
			}

			fmt.Fprintf(&bld, "\nSee documentation for more information on argument boundaries, etc.\n\n")
			return YAMLConfig{}, errors.New(bld.String())
		}

		return YAMLConfig{}, err
	}

	// Handle specific fields.
	cfg.WorkingDir, err = fixWorkingDir(cfg.WorkingDir)
	if err != nil {
		return YAMLConfig{}, err
	}

	return cfg, nil
}

func fixWorkingDir(workingDir string) (string, error) {
	// Expand '~' prefix to the user's home directory.
	if workingDir[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		workingDir = filepath.Join(homeDir, workingDir[1:])
	}

	// Use the absolute path form of the working directory.
	workingDir, err := filepath.Abs(workingDir)
	if err != nil {
		return "", err
	}

	// Verify that cfg.WorkingDir is a directory.
	finfo, err := os.Stat(workingDir)
	if err != nil {
		return "", err
	}

	if ok := finfo.IsDir(); !ok {
		return "", fmt.Errorf("invalid working directory: %s", workingDir)
	}

	return workingDir, nil
}
