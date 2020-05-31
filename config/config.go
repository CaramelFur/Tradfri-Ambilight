package config

import (
	viperMain "github.com/spf13/viper"
)

// Config holds all necessary data to stay logged in to the gateway
type Config struct {
	Address      string
	ClientID     string
	GeneratedPsk string
	StandardPsk  string

	viper *viperMain.Viper
}

// NewConfig initializes an empty config
func NewConfig() (*Config, error) {
	var viper *viperMain.Viper = viperMain.New()

	viper.AddConfigPath("$HOME/.config/")
	viper.SetConfigName("ambitradfri")
	viper.SetConfigType("yaml")

	var cfg Config = Config{
		Address:      "",
		ClientID:     "",
		GeneratedPsk: "",
		StandardPsk:  "",

		viper: viper,
	}

	err := cfg.Read()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves the current config to disk
func (c *Config) Save() error {
	c.viper.Set("login.address", c.Address)
	c.viper.Set("login.clientid", c.ClientID)
	c.viper.Set("login.psk", c.GeneratedPsk)
	c.viper.Set("psk", c.StandardPsk)

	err := c.viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

// Read reads the current config from disk into the config
func (c *Config) Read() error {
	if err := c.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viperMain.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return nil
		}

		return err
	}

	c.Address = c.viper.GetString("login.address")
	c.ClientID = c.viper.GetString("login.clientid")
	c.GeneratedPsk = c.viper.GetString("login.psk")
	c.StandardPsk = c.viper.GetString("psk")

	return nil
}
