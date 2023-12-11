package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func NewServerConfig(path, filename string) (Server, error) {
	cfgFilename := "server.yml"
	if os.Getenv("POW_CONFIG_FILE") != "" {
		cfgFilename = os.Getenv("POW_CONFIG_FILE")
	}

	cfgPath := "./config"
	if os.Getenv("POW_CONFIG_PATH") != "" {
		cfgPath = os.Getenv("POW_CONFIG_PATH")
	}

	viper.SetConfigName(cfgFilename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfgPath)

	err := viper.ReadInConfig()
	if err != nil {
		return Server{}, fmt.Errorf("Error config file: %w \n", err)
	}

	cfg := Server{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Server{}, fmt.Errorf("Unable to decode Config: %w \n", err)
	}

	return cfg, nil
}

func NewClientConfig(path, filename string) (Client, error) {
	cfgFilename := "client.yml"
	if os.Getenv("POW_CONFIG_FILE") != "" {
		cfgFilename = os.Getenv("POW_CONFIG_FILE")
	}

	cfgPath := "./config"
	if os.Getenv("POW_CONFIG_PATH") != "" {
		cfgPath = os.Getenv("POW_CONFIG_PATH")
	}

	viper.SetConfigName(cfgFilename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfgPath)

	err := viper.ReadInConfig()
	if err != nil {
		return Client{}, fmt.Errorf("Error config file: %w \n", err)
	}

	cfg := Client{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Client{}, fmt.Errorf("Unable to decode Config: %w \n", err)
	}

	return cfg, nil
}
