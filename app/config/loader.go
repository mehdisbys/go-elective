package config

import (
	"io/ioutil"

	"github.com/caarlos0/env"
	"github.com/go-playground/validator"
	"gopkg.in/yaml.v2"
)

// Envs is a struct that defines the required env variables
type Envs struct {
	Env string `env:"ENV"`
}

// Config is a struct that holds auth keys, db config, sqs, sns ... and environment variables
type Config struct {
	Env           Envs
	CountriesFile string `yaml:"countriesFile" validate:"required"`
	IntervalMs    int    `yaml:"intervalMs" validate:"required"`
	Port          int    `yaml:"port" validate:"required"`
}

// NewConfig returns a new `*Config`
func NewConfig() (*Config, error) {
	v := validator.New()

	envs, err := GetEnvConfigs(v)
	if err != nil {
		return nil, err
	}

	if envs.Env == "" {
		envs.Env = "dev"
	}

	cfg, err := GetServiceConfig(envs.Env, v)
	if err != nil {
		return nil, err
	}

	cfg.Env = *envs
	return cfg, nil
}

// GetEnvConfigs gets env variables as specified per struct at top of file
func GetEnvConfigs(v *validator.Validate) (*Envs, error) {
	c := Envs{}

	err := env.Parse(&c)
	if err != nil {
		return nil, err
	}

	err = v.Struct(c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetServiceConfig get the service config as specified per struct at top of file
func GetServiceConfig(env string, v *validator.Validate) (*Config, error) {
	source, err := ioutil.ReadFile("config" + "/" + env + ".yaml")
	if err != nil {
		return nil, err
	}

	c := Config{}

	err = yaml.Unmarshal(source, &c)
	if err != nil {
		return nil, err
	}

	err = v.StructExcept(c, "Env")
	if err != nil {
		return nil, err
	}

	return &c, nil
}
