package services

import (
	"log"
	"os"
	"path/filepath"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/spf13/viper"
)

// GetCnf get config struct
func GetCnf(c interface{}, files map[string]string) {
	
	for folder, name := range files {
		log.Printf("Loading config %s/%s", folder, name)
		err := InitFiles(folder, name)

		if err != nil {
			log.Printf("Warning: in config %s", err)
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
	viper.SetEnvPrefix("CBOTS")
	viper.SetConfigName("")
	viper.SetConfigType("env")
	viper.AddConfigPath(helpers.GetCurrDir())
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		log.Printf("Warning: %s", err)
	}
}
