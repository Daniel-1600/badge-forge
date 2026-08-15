package rules

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type Config struct {
	Rules []RuleConfig `yaml:"rules"`
}

type RuleConfig struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Enabled     bool   `yaml:"enabled"`
	Metric      string `yaml:"metric"`
	Threshold   int    `yaml:"threshold"`
	BadgeID     string `yaml:"badge_id"`
}

// LoadRules reads rule definitions from a YAML file and creates the
// corresponding executable rules using the supplied database connection.
func LoadRules(path string, db *gorm.DB) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read rules file %q: %w", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse rules file %q: %w", path, err)
	}

	loaded := make([]Rule, 0, len(config.Rules))
	for i, definition := range config.Rules {
		if !definition.Enabled {
			continue
		}

		switch definition.Type {
		case "milestone":
			if definition.Threshold <= 0 {
				return nil, fmt.Errorf("rule %q (index %d): threshold must be greater than zero", definition.Name, i)
			}

			loaded = append(loaded, &MilestoneRule{
				Threshold: definition.Threshold,
				DB:        db,
			})
		default:
			return nil, fmt.Errorf("rule %q (index %d): unknown type %q", definition.Name, i, definition.Type)
		}
	}

	return loaded, nil
}
