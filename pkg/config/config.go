package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Checks []Check `yaml:"checks"`
}

type Check struct {
	Sources  []Source `yaml:"sources"`
	Policies []Source `yaml:"policies"`
}

type Source struct {
	Path string `yaml:"path"`

	IsDir   bool `yaml:"-"`
	IsValid bool `yaml:"-"`
}

func Load(path string) (*Config, error) {
	var cfg Config

	// only accepts yaml config
	if !strings.HasSuffix(path, ".yaml") {
		return &cfg, fmt.Errorf("config file must be a .yaml file")
	}

	dir := filepath.Dir(path)

	file, err := os.ReadFile(path)
	if err != nil {
		return &cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return &cfg, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	for i, check := range cfg.Checks {
		for j, source := range check.Sources {
			if source.Path == "" {
				cfg.Checks[i].Sources[j].IsValid = false
				continue
			}
			stat, err := os.Stat(filepath.Join(dir, source.Path))
			if err != nil {
				return &cfg, fmt.Errorf("failed to stat source file: %w", err)
			}

			if stat.IsDir() {
				cfg.Checks[i].Sources[j].IsDir = true
			}
			cfg.Checks[i].Sources[j].IsValid = true
		}
		for j, policy := range check.Policies {
			stat, err := os.Stat(filepath.Join(dir, policy.Path))
			if err != nil {
				return &cfg, fmt.Errorf("failed to stat policy file %s: %w", policy.Path, err)
			}

			if stat.IsDir() {
				cfg.Checks[i].Policies[j].IsDir = true
			}
		}
	}

	return &cfg, nil
}
