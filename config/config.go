package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Name        string  `yaml:"name"`
	ContextPath string  `yaml:"contextPath"`
	TargetUrl   string  `yaml:"targetUrl"`
	Routes      []Route `yaml:"routes"`
}

type Route struct {
	Path         string   `yaml:"path"`
	Methods      []string `yaml:"methods"`
	AuthRequired bool     `yaml:"authRequired"`
}

type Config struct {
	Services []ServiceConfig `yaml:"services"`
	SECRET   string
	PORT     string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("couldnt load env")
	}
}

func LoadConfig() (*Config, error) {
	secret := os.Getenv("SECRET")
	if len(secret) <= 0 {
		return nil, errors.New("no secret found")
	}

	port := os.Getenv("PORT")
	if len(port) <= 0 {
		return nil, errors.New("no port found")
	}
	var cfg Config
	cfg.PORT = port
	cfg.SECRET = secret
	data, err := os.ReadFile("conf.yml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) GetUrlFromEndpoint(path string) (*ServiceConfig, error) {
	for _, service := range c.Services {
		if strings.HasPrefix(path, service.ContextPath) {
			return &service, nil
		}
	}

	return nil, errors.New("no matching service with that path found")
}
