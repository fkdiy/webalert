package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval int            `yaml:"interval"`
	EMail    EMailConfig    `yaml:"email"`
	Targets  []TargetConfig `yaml:"targets"`
}

type EMailConfig struct {
	From       string     `yaml:"from"`
	Recipients []string   `yaml:"recipients"`
	SMTP       SMTPConfig `yaml:"smtp"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type TargetConfig struct {
	URL      string `yaml:"url"`
	Selector string `yaml:"selector"`
}

func Load() (Config, error) {
	dat, err := os.ReadFile("webalert.config.yaml")

	if err != nil {
		return Config{}, err
	}

	conf, err := parse(dat)

	if err != nil {
		return Config{}, err
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
