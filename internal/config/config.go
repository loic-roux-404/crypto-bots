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
		err := AddCnfFile(folder, name)

		if err != nil {
			log.Printf("Warn: config issue %s", err)
		}
	}

	viper.Unmarshal(&c)
}

// AddCnfFile config from yaml
func AddCnfFile(folder string, name string) error {
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
