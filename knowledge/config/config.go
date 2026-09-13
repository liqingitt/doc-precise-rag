package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ImportProcessConfig struct {
	ImportFileTempDir *string `yaml:"import_file_temp_dir"`
}

type MineruConfig struct {
	BaseURL *string `yaml:"base_url"`
	Token   *string `yaml:"token"`
}

type VlmConfig struct {
	ModelName *string `yaml:"model_name"`
	Token     *string `yaml:"token"`
	BaseURL   *string `yaml:"base_url"`
}

type AiConfig struct {
	MineruConfig *MineruConfig `yaml:"mineru_config"`
	VlmConfig    *VlmConfig    `yaml:"vlm_config"`
}

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
	CosConfig           *CosConfig           `yaml:"cos_config"`
	MysqlConfig         *MysqlConfig         `yaml:"mysql_config"`
	AiConfig            *AiConfig            `yaml:"ai_config"`
	ImportProcessConfig *ImportProcessConfig `yaml:"import_process_config"`
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
}
