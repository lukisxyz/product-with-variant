package config

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// environment utility
// read string from env file
func loadEnvString(key string, res *string) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	*res = s
}

// read uint from env file
func loadEnvUint(key string, res *uint) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	*res = uint(num)
}

type Config struct {
	Listen   listenConfig `yaml:"listen" json:"listen"`
	DBConfig pgConfig     `yaml:"db" json:"db"`
}

func defaultConfig() Config {
	return Config{
		Listen:   defaultListenConfig(),
		DBConfig: defaultPgConfig(),
	}
}

func (c *Config) loadFromEnv() {
	c.Listen.loadFromEnv()
	c.DBConfig.loadFromEnv()
}

func loadConfigFromReader(r io.Reader, c *Config) error {
	return yaml.NewDecoder(r).Decode(c)
}

func loadConfigFromFile(fName string, c *Config) error {
	_, err := os.Stat(fName)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath.Clean(fName))
	if err != nil {
		return err
	}

	defer f.Close()
	return loadConfigFromReader(f, c)
}

func LoadConfig(fileName string) Config {
	// load from default configuration
	cfg := defaultConfig()

	// load from environment
	cfg.loadFromEnv()

	err := loadConfigFromFile(fileName, &cfg)
	if err != nil {
		return cfg
	}
	return cfg
}
