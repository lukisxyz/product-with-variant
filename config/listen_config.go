package config

import "fmt"

type listenConfig struct {
	Host    string `yaml:"host" json:"host"`
	Port    uint   `yaml:"port" json:"port"`
	ReadTO  uint   `yaml:"read" json:"read"`
	WriteTO uint   `yaml:"write" json:"write"`
	IdleTO  uint   `yaml:"idle" json:"idle"`
}

func (l listenConfig) Address() string {
	return fmt.Sprintf(
		"%s:%d",
		l.Host,
		l.Port,
	)
}

func defaultListenConfig() listenConfig {
	return listenConfig{
		Host: "127.0.0.1",
		Port: 8080,
	}
}

func (l *listenConfig) loadFromEnv() {
	loadEnvString("LISTEN_HOST", &l.Host)
	loadEnvUint("LISTEN_PORT", &l.Port)
	loadEnvUint("READ_TIMEOUT", &l.ReadTO)
	loadEnvUint("WRITE_TIMEOUT", &l.WriteTO)
	loadEnvUint("IDLE_TIMEOUT", &l.IdleTO)
}
