package conf

import (
	"os"
	"profile_service/pkg/models"
	"strconv"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

// Конфиги приложения
type Config struct {
	Server struct {
		Port             string `yaml:"port"`
		Host             string `yaml:"host"`
		Prefix           string `yaml:"prefix"`
		JWTAccessSecret  string `yaml:"jwt_access_secret"`
		JWTRefreshSecret string `yaml:"jwt_refresh_secret"`
	} `yaml:"server"`
	AuthService struct {
		Port            string `yaml:"port"`
		Host            string `yaml:"host"`
		ProfileDetails  string `yaml:"profile_details"`
		TokenValidation string `yaml:"token_validation"`
		Retries         int    `yaml:"retries"`
		RetryDelay      int    `yaml:"retry_delay"`
	} `yaml:"auth_service"`
	Database struct {
		Username string                     `yaml:"user"`
		Password string                     `yaml:"password"`
		Users    map[uuid.UUID]*models.User `yaml:"users"`
	} `yaml:"database"`
}

func New() *Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func (c *Config) GetAuthServiceAddr() string {
	return c.AuthService.Host + ":" + c.AuthService.Port
}

func (c *Config) GetAuthServiceProfileDetailsUrl() string {
	return c.GetAuthServiceAddr() + "/" + c.AuthService.ProfileDetails
}

func (c *Config) GetAuthServiceTokenValidationUrl() string {
	return c.GetAuthServiceAddr() + "/" + c.AuthService.TokenValidation
}

func (c *Config) AuthServiceRetries() int {
	return c.AuthService.Retries
}

func (c *Config) AuthServiceRetryDelay() int {
	return c.AuthService.RetryDelay
}

func (c *Config) ServerAddr() string {
	return c.Server.Host + ":" + c.Server.Port
}

func getenvInt(key string) int {
	v, err := strconv.Atoi(key)
	if err != nil {
		return 0
	}
	return v
}
