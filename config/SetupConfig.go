package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func Create() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			os.Create("config.json")
		} else {
			log.Fatal("Config not set")
		}
	}

	if !viper.GetBool("setup_config") {
		var laravelPath string
		var golangPath string
		var rustPath string

		viper.Set("setup_config", true)

		fmt.Print("Laravel Path:[default:/var/www/html/]")
		fmt.Scanln(&laravelPath)
		if laravelPath != "" {
			viper.Set("laravel_path", laravelPath)
		}

		fmt.Print("Golang Path:[default:~/go/src/]")
		fmt.Scanln(&golangPath)
		if golangPath != "" {
			viper.Set("golang_path", golangPath)
		}

		fmt.Print("RUST Path:[default:~/Documents/Code/Rust/]")
		fmt.Scanln(&rustPath)
		if rustPath != "" {
			viper.Set("rust_path", rustPath)
		}

		viper.WriteConfig()
	}
}
