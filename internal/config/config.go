package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Get get config struct
func Get(c interface{}, files map[string]string) {
	for folder, name := range files {
		log.Printf("Info: Loading config %s/%s", folder, name)
		err := InitFiles(folder, name)

		if err != nil {
			log.Printf("Warn: config issue %s", err)
		}
	}

	InitEnv()

	viper.Unmarshal(&c)
}

// InitFiles config from yaml
func InitFiles(folder string, name string) error {
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	wd, _ := os.Getwd()
	viper.AddConfigPath(filepath.Join(wd, "config", folder))
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	return nil
}

// InitEnv parse secured env vars
func InitEnv() {
	wd, _ := os.Getwd()
	viper.SetEnvPrefix("cbots")
	viper.AddConfigPath(wd)
	viper.SetConfigName("")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		log.Printf("Warning: %s", err)
	}
}
