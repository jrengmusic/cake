package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the cake configuration
type Config struct {
	AutoScan   AutoScanConfig   `toml:"auto_scan"`
	Appearance AppearanceConfig `toml:"appearance"`
	Build      BuildConfig      `toml:"build"`
}

// BuildConfig holds build-related settings (last chosen options)
type BuildConfig struct {
	LastProject       string `toml:"last_project"`
	LastConfiguration string `toml:"last_configuration"`
}

// AutoScanConfig holds auto-scan settings
type AutoScanConfig struct {
	Enabled         bool `toml:"enabled"`
	IntervalMinutes int  `toml:"interval_minutes"`
}

// AppearanceConfig holds appearance settings
type AppearanceConfig struct {
	Theme string `toml:"theme"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		AutoScan: AutoScanConfig{
			Enabled:         true,
			IntervalMinutes: 10,
		},
		Appearance: AppearanceConfig{
			Theme: "gfx",
		},
		Build: BuildConfig{
			LastProject:       "",
			LastConfiguration: "Debug",
		},
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".cake/config.toml"
	}
	return filepath.Join(home, ".config", "cake", "config.toml")
}

// GetThemesPath returns the path to the themes directory
func GetThemesPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".cake/themes"
	}
	return filepath.Join(home, ".config", "cake", "themes")
}

// Load loads the configuration from the config file
// If the config file doesn't exist, it creates a default config
func Load() (*Config, error) {
	configPath := GetConfigPath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		cfg := DefaultConfig()
		if err := Save(cfg); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse TOML
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for missing values
	if cfg.AutoScan.IntervalMinutes == 0 {
		cfg.AutoScan.IntervalMinutes = 10
	}
	if cfg.Appearance.Theme == "" {
		cfg.Appearance.Theme = "gfx"
	}

	return &cfg, nil
}

// Save saves the configuration to the config file
func Save(cfg *Config) error {
	configPath := GetConfigPath()

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to TOML
	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Theme returns the current theme name
func (c *Config) Theme() string {
	return c.Appearance.Theme
}

// SetTheme updates the theme and saves the config
func (c *Config) SetTheme(theme string) error {
	c.Appearance.Theme = theme
	return Save(c)
}

// IsAutoScanEnabled returns whether auto-scan is enabled
func (c *Config) IsAutoScanEnabled() bool {
	return c.AutoScan.Enabled
}

// SetAutoScanEnabled updates auto-scan enabled status and saves
func (c *Config) SetAutoScanEnabled(enabled bool) error {
	c.AutoScan.Enabled = enabled
	return Save(c)
}

// AutoScanInterval returns the auto-scan interval in minutes
func (c *Config) AutoScanInterval() int {
	return c.AutoScan.IntervalMinutes
}

// SetAutoScanInterval updates the auto-scan interval and saves
func (c *Config) SetAutoScanInterval(minutes int) error {
	c.AutoScan.IntervalMinutes = minutes
	return Save(c)
}

// LastProject returns the last chosen project
func (c *Config) LastProject() string {
	return c.Build.LastProject
}

// SetLastProject updates the last chosen project and saves
func (c *Config) SetLastProject(project string) error {
	c.Build.LastProject = project
	return Save(c)
}

// LastConfiguration returns the last chosen configuration (Debug/Release)
func (c *Config) LastConfiguration() string {
	if c.Build.LastConfiguration == "" {
		return "Debug"
	}
	return c.Build.LastConfiguration
}

// SetLastConfiguration updates the last chosen configuration and saves
func (c *Config) SetLastConfiguration(configuration string) error {
	c.Build.LastConfiguration = configuration
	return Save(c)
}
