package config

import "fmt"

type pgConfig struct {
	Host string `yaml:"host" json:"host"`
	Port uint   `yaml:"port" json:"port"`

	DBName  string `yaml:"db_name" json:"db_name"`
	SslMode string `yaml:"ssl_mode" json:"ssl_mode"`

	Username string `yaml:"username" json:"username"`
	Secret   string `yaml:"secret" json:"secret"`
}

func (p *pgConfig) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s database=%s sslmode=%s",
		p.Host,
		p.Port,
		p.Username,
		p.Secret,
		p.DBName,
		p.SslMode,
	)
}

func defaultPgConfig() pgConfig {
	return pgConfig{
		Host:     "localhost",
		Port:     5432,
		DBName:   "product",
		SslMode:  "disable",
		Username: "postgres",
		Secret:   "",
	}
}

func (p *pgConfig) loadFromEnv() {
	loadEnvString("DB_HOST", &p.Host)
	loadEnvUint("DB_PORT", &p.Port)
	loadEnvString("DB_NAME", &p.DBName)
	loadEnvString("DB_SSL", &p.SslMode)
	loadEnvString("DB_USER", &p.Username)
	loadEnvString("DB_SECRET", &p.Secret)
}
