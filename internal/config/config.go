package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval int            `yaml:"interval"`
	Jitter   int            `yaml:"jitter"`
	EMail    EMailConfig    `yaml:"email"`
	Targets  []TargetConfig `yaml:"targets"`
}

type EMailConfig struct {
	Mode       EMailMode  `yaml:"mode"`
	From       string     `yaml:"from"`
	Recipients []string   `yaml:"recipients"`
	SMTP       SMTPConfig `yaml:"smtp"`
}

type EMailMode string

const (
	EMailModePerChange EMailMode = "per_change"
	EMailModeDigest    EMailMode = "digest"
)

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type TargetConfig struct {
	URL      string `yaml:"url"`
	Selector string `yaml:"selector"`
}

func Load(path string) (Config, error) {
	dat, err := os.ReadFile(path)

	if err != nil {
		return Config{}, err
	}

	conf, err := parse(dat)

	if err != nil {
		return Config{}, err
	}

	if conf.Jitter < 0 {
		return Config{}, fmt.Errorf("jitter must not be negative")
	}

	if conf.Jitter >= conf.Interval {
		return Config{}, fmt.Errorf("jitter must be smaller than interval")
	}

	return conf, nil
}

func parse(dat []byte) (Config, error) {
	var conf Config

	err := yaml.Unmarshal(dat, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
