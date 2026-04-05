package config

import (
	"testing"

	"github.com/jrengmusic/cake/internal"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	if cfg.Appearance.Theme != "gfx" {
		t.Errorf("expected Theme %q, got %q", "gfx", cfg.Appearance.Theme)
	}
	if !cfg.AutoScan.Enabled {
		t.Error("expected AutoScan.Enabled to be true")
	}
	if cfg.AutoScan.IntervalMinutes != internal.DefaultAutoScanInterval {
		t.Errorf("expected AutoScan.IntervalMinutes %d, got %d", internal.DefaultAutoScanInterval, cfg.AutoScan.IntervalMinutes)
	}
	if cfg.Build.LastProject != "" {
		t.Errorf("expected Build.LastProject empty, got %q", cfg.Build.LastProject)
	}
	if cfg.Build.LastConfiguration != internal.ConfigDebug {
		t.Errorf("expected Build.LastConfiguration %q, got %q", internal.ConfigDebug, cfg.Build.LastConfiguration)
	}
}

func TestTheme(t *testing.T) {
	tests := []struct {
		name     string
		theme    string
		expected string
	}{
		{"gfx theme", "gfx", "gfx"},
		{"dark theme", "dark", "dark"},
		{"empty theme", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Appearance: AppearanceConfig{Theme: tt.theme},
			}
			got := cfg.Theme()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestIsAutoScanEnabled(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{"enabled true", true, true},
		{"enabled false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				AutoScan: AutoScanConfig{Enabled: tt.enabled},
			}
			got := cfg.IsAutoScanEnabled()
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestAutoScanInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval int
		expected int
	}{
		{"default interval", internal.DefaultAutoScanInterval, internal.DefaultAutoScanInterval},
		{"custom interval", 30, 30},
		{"zero interval", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				AutoScan: AutoScanConfig{IntervalMinutes: tt.interval},
			}
			got := cfg.AutoScanInterval()
			if got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestLastProject(t *testing.T) {
	tests := []struct {
		name        string
		lastProject string
		expected    string
	}{
		{"non-empty project", "MyProject", "MyProject"},
		{"empty project", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Build: BuildConfig{LastProject: tt.lastProject},
			}
			got := cfg.LastProject()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestLastConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		configuration string
		expected      string
	}{
		{"debug configuration", internal.ConfigDebug, internal.ConfigDebug},
		{"release configuration", internal.ConfigRelease, internal.ConfigRelease},
		{"empty falls back to Debug", "", internal.ConfigDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Build: BuildConfig{LastConfiguration: tt.configuration},
			}
			got := cfg.LastConfiguration()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
