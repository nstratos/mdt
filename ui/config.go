package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type ConfigField string

func (cf ConfigField) Val() string {
	return string(cf)
}

const (
	ConfigMode      ConfigField = "Mode"
	ConfigTotalTime ConfigField = "TotalTime"
	ConfigOffset    ConfigField = "Offset"
	ConfigBaseHz    ConfigField = "BaseHz"
	ConfigStartHz   ConfigField = "StartHz"
	ConfigEndHz     ConfigField = "EndHz"
)

type Config struct {
	Mode      rune
	TotalTime int
	Offset    int
	BaseHz    float64
	StartHz   float64
	EndHz     float64
}

func (c Config) Save() error {
	return writeConfig(c)
}

func writeConfig(c Config) error {
	json, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", json, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Load() error {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		c := Config{'A', 30, 5, 100, 15.00, 8.00}
		if err := writeConfig(c); err != nil {
			return err
		}
		b, err = ioutil.ReadFile("config.json")
		if err != nil {
			return err
		}
	}
	//s, _ := strconv.Unquote(string(b))
	//err = json.Unmarshal([]byte(s), c)
	err = json.Unmarshal(b, c)
	if err != nil {
		return fmt.Errorf("loading config error: %v", err)
	}
	return nil
}

func (c *Config) Update(m map[string]interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp, c)
	if err != nil {
		return fmt.Errorf("Unmarshal: %v\n", err)
	}
	return nil
}

func GetConfig() Config {
	return config
}

func UpdateConfig(c Config) {
	config = c
}

func (c Config) ModeS() string {
	return strconv.QuoteRuneToASCII(c.Mode)
}

func (c Config) TotalTimeS() string {
	return fmt.Sprintf("%v min", c.TotalTime)
}

func (c Config) OffsetS() string {
	return fmt.Sprintf("%v min", c.Offset)
}

func (c Config) BaseHzS() string {
	return fmt.Sprintf("%.2f hz", c.BaseHz)
}

func (c Config) StartHzS() string {
	return fmt.Sprintf("%.2f hz", c.StartHz)
}

func (c Config) EndHzS() string {
	return fmt.Sprintf("%.2f hz", c.EndHz)
}
