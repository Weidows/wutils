package config

import (
	"os"
	"testing"
)

func TestConfigVersionDefault(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ConfigVersion != CurrentConfigVersion {
		t.Errorf("expected version %d, got %d", CurrentConfigVersion, cfg.ConfigVersion)
	}
}

func TestConfigVersionMigration(t *testing.T) {
	// Simulate an old config (version 0)
	cfg := Config{
		ConfigVersion: 0,
		App: AppConfig{
			Name:    "old-wutils",
			Version: "0.5.0",
		},
		Refresh: 10,
		Cmd: CmdConfig{
			Dsg: DSGConfig{
				Parallel: true,
				Disk:     []string{"E:"},
				Delay:    30,
			},
			Ol: OLConfig{
				Parallel: true,
				Delay:    2,
			},
			Buffer: BufferConfig{
				Strategy: "balanced",
			},
		},
	}

	// Run migrations
	RunMigrations(&cfg)

	// Should be migrated to version 1
	if cfg.ConfigVersion != 1 {
		t.Errorf("expected version 1 after migration, got %d", cfg.ConfigVersion)
	}
}

func TestConfigVersionAlreadyCurrent(t *testing.T) {
	cfg := DefaultConfig()

	// Run migrations on already-current config should be a no-op
	RunMigrations(&cfg)

	if cfg.ConfigVersion != CurrentConfigVersion {
		t.Errorf("expected unchanged version %d, got %d", CurrentConfigVersion, cfg.ConfigVersion)
	}
}

func TestLocaleDefault(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Locale != "auto" {
		t.Errorf("expected 'auto', got %s", cfg.Locale)
	}
}

func TestLoadWithMigration(t *testing.T) {
	// Create an old-style config (without config_version)
	oldYAML := `
app:
  name: legacy-wutils
  version: "0.9.0"
refresh: 15
cmd:
  dsg:
    parallel: true
    disk:
      - "D:"
    delay: 60
  ol:
    parallel: true
    delay: 3
  buffer:
    strategy: download
`
	dir := t.TempDir()
	path := dir + "/old_config.yml"
	if err := os.WriteFile(path, []byte(oldYAML), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Load should migrate automatically
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v\nYAML:\n%s", err, oldYAML)
	}

	// Config should be migrated to current version
	if cfg.ConfigVersion != CurrentConfigVersion {
		t.Errorf("expected version %d after load+migrate, got %d", CurrentConfigVersion, cfg.ConfigVersion)
	}

	// Values from the file should be preserved
	if cfg.App.Name != "legacy-wutils" {
		t.Errorf("expected 'legacy-wutils', got %s", cfg.App.Name)
	}
	if cfg.Cmd.Dsg.Delay != 60 {
		t.Errorf("expected 60, got %d", cfg.Cmd.Dsg.Delay)
	}

	// Defaults should be filled
	if cfg.Locale != "auto" {
		t.Errorf("expected 'auto', got %s", cfg.Locale)
	}
}
