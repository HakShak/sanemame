package config

import (
	"log"
	"os"
)

import "github.com/spf13/viper"

const MameRepo = "mame.repo"
const GithubReleaseApi = "github.api.release"

func SetupConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	//Check if we have a config file first
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if _, ok := err.(*os.PathError); ok {
		pwd, _ := os.Getwd()
		log.Println("No config found in %s", pwd)
	}

	viper.SetDefault(MameRepo, "mamedev/mame")
	viper.SetDefault(GithubReleaseApi, "https://api.github.com/repos/%s/releases")
}