package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
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
		return err
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() Config {
	return config
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
