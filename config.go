package main

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set.
func initConfig(configPath string) error {
	if configPath != "" {
		// Use config file from flag
		viper.SetConfigFile(configPath)
	} else {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf("Cannot find home directory: %s", err)
		}

		// Search config in home directory with name ".cyborg" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".cyborg")
	}

	viper.SetEnvPrefix("cyborg")
	viper.AutomaticEnv()

	viper.ReadInConfig()
	return nil
}
