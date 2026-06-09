package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	CurrentContext   string `mapstructure:"current_context"`
	DefaultItemState string `mapstructure:"default_state"`
	v                *viper.Viper
}

const (
	DefaultItemState = "todo"
)

func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ctx")
}

func Path() string {
	return filepath.Join(Dir(), "config.yaml")
}

func Load() (*Config, error) {
	if err := os.MkdirAll(Dir(), 0755); err != nil {
		return nil, fmt.Errorf("could not create config dir: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(Path())
	v.SetConfigType("yaml")

	if _, err := os.Stat(Path()); os.IsNotExist(err) {
		v.Set("current_context", "")
		v.Set("default_state", DefaultItemState)
		if err := v.WriteConfigAs(Path()); err != nil {
			return nil, fmt.Errorf("could not write config: %w", err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	cfg := &Config{v: v}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}
	return cfg, nil
}

func (c *Config) Save() error {
	c.v.Set("current_context", c.CurrentContext)
	c.v.Set("default_state", c.DefaultItemState)
	if err := c.v.WriteConfig(); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	return nil
}

// SetCurrentContext sets the current context and persists the config.
func (c *Config) SetCurrentContext(name string) error {
	c.CurrentContext = name
	return c.Save()
}

// ClearCurrentContext clears the current context (sets it to global) and persists the config.
func (c *Config) ClearCurrentContext() error {
	c.CurrentContext = ""
	return c.Save()
}

// GetContextDir returns the filesystem directory for a context's stored files.
func GetContextDir(name string) string {
	return filepath.Join(Dir(), "contexts", name)
}
