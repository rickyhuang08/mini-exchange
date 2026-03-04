package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Server struct {
	AppEnv         string
	Port           string `yaml:"port"`
	ConfigYamlPath string
}

type Finnhub struct {
	Scheme  string `yaml:"scheme"`
	Host    string `yaml:"host"`
	ApiKey  string
	Symbols []string `yaml:"symbols"`
}

type JWT struct {
	PrivateKeyPath string
	PublicKeyPath  string
	Expiration     int `yaml:"expiration"`
}

type Logger struct {
	AccessPath string `yaml:"access_path"`
	ErrorPath  string `yaml:"error_path"`
}

type Config struct {
	Server  Server  `yaml:"server"`
	Finnhub Finnhub `yaml:"finnhub"`
	JWT     JWT     `yaml:"jwt"`
	Logger  Logger  `yaml:"logger"`
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = ".env." + env
	}

	var cfg Config
	cfg.Server.AppEnv = env
	cfg, err := LoadConfigFromFile(cfg, envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from file %s: %w", envPath, err)
	}

	configYamlFile := cfg.Server.ConfigYamlPath
	if configYamlFile == "" {
		configYamlFile = "config/config." + env + ".yaml"
	}

	file, err := os.ReadFile(configYamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config YAML file %s: %w", configYamlFile, err)
	}

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config YAML file %s: %w", configYamlFile, err)
	}

	log.Printf("Configuration loaded successfully from %s and %s", envPath, configYamlFile)
	return &cfg, nil
}

func LoadConfigFromFile(cfg Config, path string) (Config, error) {
	// Load the corresponding config file based on the environment
	if err := godotenv.Load(path); err != nil {
		return Config{}, fmt.Errorf("error loading %s file", path)
	}

	log.Printf("Loaded environment variables from %s", path)

	if configYamlPath := os.Getenv("CONFIG_YAML_PATH"); configYamlPath != "" {
		cfg.Server.ConfigYamlPath = configYamlPath
		log.Printf("Using config YAML path from environment variable: %s", configYamlPath)
	}

	if configFinnhubApiKey := os.Getenv("FINNHUB_API_KEY"); configFinnhubApiKey != "" {
		cfg.Finnhub.ApiKey = configFinnhubApiKey
	}

	if privateKeyPath := os.Getenv("PRIVATE_KEY_PATH"); privateKeyPath != "" {
		cfg.JWT.PrivateKeyPath = privateKeyPath
		log.Printf("Using JWT private key path from environment variable: %s", privateKeyPath)
	}

	if publicKeyPath := os.Getenv("PUBLIC_KEY_PATH"); publicKeyPath != "" {
		cfg.JWT.PublicKeyPath = publicKeyPath
		log.Printf("Using JWT public key path from environment variable: %s", publicKeyPath)
	}

	return cfg, nil
}
