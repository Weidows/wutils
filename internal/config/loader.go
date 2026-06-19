package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Load reads a YAML config file, merges with defaults, and validates.
// If path is empty, returns DefaultConfig.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig() // start with compiled-in defaults

	if path == "" {
		return &cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %w", path, err)
	}

	// Overlay file values onto defaults
	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return nil, fmt.Errorf("failed to parse config %s: %w", path, err)
	}
	cfg.Merge(&fileCfg)

	// Run migrations if the config is from an older version
	RunMigrations(&cfg)

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Save back if migration ran (upgrades the file)
	if cfg.ConfigVersion > fileCfg.ConfigVersion {
		if err := Save(&cfg, path); err != nil {
			// Non-fatal: log but continue
		}
	}

	return &cfg, nil
}

// Save writes the Config to a YAML file at the given path.
func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// Template returns the default config template YAML bytes
// from the executable's config/ directory.
func Template() ([]byte, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(filepath.Dir(execPath), "config", "wutils.yml")
	return os.ReadFile(configPath)
}

// EnsureUserConfig creates a default user config at the given path
// if it does not already exist.
func EnsureUserConfig(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil // already exists
	}

	templateData, err := Template()
	if err != nil {
		// If template not found, write a minimal default
		cfg := DefaultConfig()
		return Save(&cfg, path)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return os.WriteFile(path, templateData, 0644)
}

// MergeYAMLNodes recursively merges a template YAML node tree with a user YAML node tree.
// User values take precedence for existing keys; new keys from template are added.
// This preserves YAML comments in the user config.
func MergeYAMLNodes(template, user *yaml.Node) *yaml.Node {
	if template == nil {
		return user
	}
	if user == nil {
		return template
	}

	result := *template

	if template.Kind == yaml.MappingNode && user.Kind == yaml.MappingNode {
		templateMap := nodeToMap(template)
		userMap := nodeToMap(user)

		mergedContent := make([]*yaml.Node, 0)

		for i := 0; i < len(template.Content); i += 2 {
			key := template.Content[i]
			templateValue := template.Content[i+1]

			keyStr := key.Value
			if userValue, exists := userMap[keyStr]; exists {
				mergedContent = append(mergedContent, key)
				mergedContent = append(mergedContent, MergeYAMLNodes(templateValue, userValue))
			} else {
				mergedContent = append(mergedContent, key)
				mergedContent = append(mergedContent, templateValue)
			}
		}

		for i := 0; i < len(user.Content); i += 2 {
			key := user.Content[i]
			keyStr := key.Value
			if _, exists := templateMap[keyStr]; !exists {
				mergedContent = append(mergedContent, key)
				mergedContent = append(mergedContent, user.Content[i+1])
			}
		}

		result.Content = mergedContent
	} else if user.Kind != 0 {
		return user
	}

	return &result
}

func nodeToMap(node *yaml.Node) map[string]*yaml.Node {
	result := make(map[string]*yaml.Node)
	if node.Kind != yaml.MappingNode {
		return result
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]
		if key.Kind == yaml.ScalarNode && key.Value != "" {
			result[key.Value] = value
		}
	}
	return result
}
