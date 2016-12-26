package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const maxHz = 999.99

// ConfigField is a key of each configuration value. Each input holds a
// configuration field so it's easier to update the config.
type ConfigField string

// Val returns the configuration field's value.
func (cf ConfigField) Val() string {
	return string(cf)
}

const (
	configFolder = ".mdt"
	configFile   = "config.json"

	configMode      ConfigField = "Mode"
	configTotalTime ConfigField = "TotalTime"
	configOffset    ConfigField = "Offset"
	configBaseHz    ConfigField = "BaseHz"
	configStartHz   ConfigField = "StartHz"
	configEndHz     ConfigField = "EndHz"
)

var defaultConfig = Config{"Binaural", 30, 5, 100, 15.00, 8.00}

// Config represents the program's configuration.
type Config struct {
	Mode      string // A = Binaural, B = Isochronic
	TotalTime int
	Offset    int
	BaseHz    float64
	StartHz   float64
	EndHz     float64
}

// Validate returns an error if the values of the configuration are not valid.
func (c Config) Validate() error {
	if c.Offset >= c.TotalTime {
		return errors.New("Offset must be lower than total time")
	}
	if c.BaseHz > maxHz || c.StartHz > maxHz || c.EndHz > maxHz {
		return errors.New("Hz value way too high")
	}
	return nil
}

// Save writes the configuration to config.json file.
func (c Config) Save() error {
	return writeConfig(c)
}

func writeConfig(c Config) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	u, err := user.Current()
	if err != nil {
		return err
	}
	if err := createConfigFolderIfNotExist(); err != nil {
		return err
	}
	p := filepath.Join(u.HomeDir, configFolder, configFile)
	return ioutil.WriteFile(p, b, 0644)
}

func createConfigFolderIfNotExist() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	path := filepath.Join(u.HomeDir, configFolder)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// Load loads configuration from the config.json file.
func (c *Config) Load() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	configPath := filepath.Join(u.HomeDir, configFolder, configFile)
	// If config does not exist, it gets created with default values.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := writeConfig(defaultConfig); err != nil {
			return err
		}
	}
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	//s, _ := strconv.Unquote(string(b))
	//err = json.Unmarshal([]byte(s), c)
	err = json.Unmarshal(b, c)
	if err != nil {
		return fmt.Errorf("loading config error: %v", err)
	}
	return nil
}

// Update updates the configuration values be accepting a map of these values
// using the key of each configuration field.
func (c *Config) Update(m map[string]interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp, c)
}

// GetConfig returns the configuration.
func GetConfig() Config {
	return config
}

// UpdateConfig updates the configuration.
func UpdateConfig(c Config) {
	config = c
}

// ModeS returns a string representation of the mode.
func (c Config) ModeS() string {
	return c.Mode
}

// TotalTimeS returns a string representation of the total time.
func (c Config) TotalTimeS() string {
	return fmt.Sprintf("%v min", c.TotalTime)
}

// OffsetS returns a string representation of the offset.
func (c Config) OffsetS() string {
	return fmt.Sprintf("%v min", c.Offset)
}

// BaseHzS returns a string representation of the base hz.
func (c Config) BaseHzS() string {
	return fmt.Sprintf("%.2f hz", c.BaseHz)
}

// StartHzS returns a string representation of the start hz.
func (c Config) StartHzS() string {
	return fmt.Sprintf("%.2f hz", c.StartHz)
}

// EndHzS returns a string representation of the end hz.
func (c Config) EndHzS() string {
	return fmt.Sprintf("%.2f hz", c.EndHz)
}
