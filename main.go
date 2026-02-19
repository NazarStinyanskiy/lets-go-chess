package main

import (
	"lets-go-chess/cli"
	"lets-go-chess/server"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	loadConfig()
	chooseMode(1)
}

func loadConfig() {
	env := os.Getenv("ENV")
	viper.SetConfigName("config-" + env)
	viper.AddConfigPath("./cfg/")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func chooseMode(mode int) {
	switch mode {
	case 0:
		cli.StartGame()
	case 1:
		server.StartServer()
	}
}
