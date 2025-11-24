package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbName     string `mapstructure:"DB_NAME"`
	Port       string `mapstructure:"PORT"`
}

func Load() (config Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../..")
	viper.AddConfigPath("../../config")
	viper.AddConfigPath("../../..")
	viper.AddConfigPath("../../../config")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found (OK, will use env vars)")
		} else {
			log.Println("Error reading config file")
		}
	}

	viper.AutomaticEnv()

	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_PORT")
	_ = viper.BindEnv("DB_USER")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("PORT")

	log.Printf("Config loaded: %#v", viper.AllSettings())

	err = viper.Unmarshal(&config)

	return
}
