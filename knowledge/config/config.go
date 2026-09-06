package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type CosConfig struct {
	Bucket    *string `yaml:"bucket"`
	Domain    *string `yaml:"domain"`
	Scheme    *string `yaml:"scheme"`
	SecretId  *string `yaml:"secret_id"`
	SecretKey *string `yaml:"secret_key"`
}

type MysqlConfig struct {
	Host     *string `yaml:"host"`
	Port     *int64  `yaml:"port"`
	User     *string `yaml:"user"`
	Password *string `yaml:"password"`
	DBName   *string `yaml:"db_name"`
}

type Config struct {
	CosConfig   *CosConfig   `yaml:"cos_config"`
	MysqlConfig *MysqlConfig `yaml:"mysql_config"`
}

var AppConfig *Config

func init() {
	configRaw, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configRaw, &AppConfig)
	if err != nil {
		panic(err)
	}
	var cfg Config
	err = yaml.Unmarshal(configRaw, &cfg)

	if err != nil {
		panic(err)
	}
	AppConfig = &cfg
}
