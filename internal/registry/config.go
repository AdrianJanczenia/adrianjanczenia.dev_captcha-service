package registry

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		HTTPPort string `yaml:"httpPort"`
	} `yaml:"server"`
	Infrastructure struct {
		Retry struct {
			MaxAttempts  int           `yaml:"maxAttempts"`
			DelaySeconds time.Duration `yaml:"delaySeconds"`
		} `yaml:"retry"`
	} `yaml:"infrastructure"`
	Redis struct {
		URL string `yaml:"url"`
	} `yaml:"redis"`
	Security struct {
		HmacSecret string `yaml:"hmacSecret"`
		Difficulty int    `yaml:"difficulty"`
		TtlMinutes int    `yaml:"ttlMinutes"`
	} `yaml:"security"`
	Captcha struct {
		TtlMinutes int `yaml:"ttlMinutes"`
		MaxTries   int `yaml:"maxTries"`
	} `yaml:"captcha"`
}

var Cfg *Config

func LoadConfig() (*Config, error) {
	type yamlConfig struct {
		Server struct {
			HTTPPort string `yaml:"httpPort"`
		} `yaml:"server"`
		Infrastructure struct {
			Retry struct {
				MaxAttempts  int `yaml:"maxAttempts"`
				DelaySeconds int `yaml:"delaySeconds"`
			} `yaml:"retry"`
		} `yaml:"infrastructure"`
		Redis struct {
			URL string `yaml:"url"`
		} `yaml:"redis"`
		Security struct {
			HmacSecret string `yaml:"hmacSecret"`
			Difficulty int    `yaml:"difficulty"`
			TtlMinutes int    `yaml:"ttlMinutes"`
		} `yaml:"security"`
		Captcha struct {
			TtlMinutes int `yaml:"ttlMinutes"`
			MaxTries   int `yaml:"maxTries"`
		} `yaml:"captcha"`
	}

	env := os.Getenv("APP_ENV")
	if env != "production" {
		env = "local"
	}
	configPath := filepath.Join("config", env, "config.yml")
	log.Printf("INFO: loading configuration from %s", configPath)

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var yc yamlConfig
	if err := yaml.NewDecoder(f).Decode(&yc); err != nil {
		return nil, err
	}

	cfg := &Config{}
	cfg.Server.HTTPPort = yc.Server.HTTPPort
	cfg.Infrastructure.Retry.MaxAttempts = yc.Infrastructure.Retry.MaxAttempts
	cfg.Infrastructure.Retry.DelaySeconds = time.Duration(yc.Infrastructure.Retry.DelaySeconds) * time.Second
	cfg.Redis.URL = yc.Redis.URL
	cfg.Security.HmacSecret = yc.Security.HmacSecret
	cfg.Security.Difficulty = yc.Security.Difficulty
	cfg.Security.TtlMinutes = yc.Security.TtlMinutes
	cfg.Captcha.TtlMinutes = yc.Captcha.TtlMinutes
	cfg.Captcha.MaxTries = yc.Captcha.MaxTries

	overrideFromEnv("REDIS_URL", &cfg.Redis.URL)
	overrideFromEnv("HMAC_SECRET", &cfg.Security.HmacSecret)

	return cfg, nil
}

func overrideFromEnv(envKey string, configValue *string) {
	if value, exists := os.LookupEnv(envKey); exists && value != "" {
		*configValue = value
	}
}
