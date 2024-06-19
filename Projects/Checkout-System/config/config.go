package config

import (
	"fmt"

	"Checkout-System/helper/json"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type (
	ServerConfig struct {
		Host string `validate:"required"`
		Port int    `validate:"required,gt=0"`
	}

	MongoConfig struct {
		Host     string `validate:"required"`
		AppName  string `validate:"required"`
		Database string `validate:"required"`
		UserName string `validate:"required"`
		Password string `validate:"required"`
	}

	Config struct {
		Server ServerConfig
		Mongo  MongoConfig
	}
)

var k = koanf.New(".")

func LoadConfig() (*Config, error) {
	if err := k.Load(file.Provider("./config/config.yaml"), yaml.Parser()); err != nil {
		return nil, err
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	fmt.Println(string(json.Parse(config)))

	return &config, nil
}
