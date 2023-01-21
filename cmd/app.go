package cmd

import (
	"fmt"
	"log"

	"github.com/projectCreator/config"
	creator "github.com/projectCreator/lib/projectCreator"
	"github.com/projectCreator/lib/utils"
	"github.com/projectCreator/lib/virtualHost"
	"github.com/spf13/viper"
)

func Create() {
	var projectName string
	var projectType string

	config.Create()

	fmt.Print("Please enter project name: ")
	_, err := fmt.Scan(&projectName)
	if err != nil {
		fmt.Errorf("Project name is incorrect %v", err)
	}

	fmt.Println("Please select your project:")
	fmt.Println("1: create laravel project")
	fmt.Println("2: create golang project")
	fmt.Println("3: create rust project")

	fmt.Print(">> ")
	n, err := fmt.Scan(&projectType)

	if err != nil {
		fmt.Errorf("Wrong select Project %v %v", err, n)
	}
	switch projectType {
	case "1":
		creator.LaravelCreator(projectName)
		virtualHost.Create(fmt.Sprintf("%s%s", viper.GetString("laravel_path"), projectName), projectName)
		utils.OpenCodeinVscode(fmt.Sprintf("%s%s", viper.GetString("laravel_path"), projectName))
		utils.LinkGenerator(projectName)
		break
	case "2":
		creator.GolangCreator(projectName)
		utils.OpenCodeinVscode(fmt.Sprintf("%s%s", viper.GetString("golang_path"), projectName))
		break
	case "3":
		creator.RustCreator(projectName)
		utils.OpenCodeinVscode(fmt.Sprintf("%s%s", viper.GetString("rust_path"), projectName))
		break
	default:
		log.Fatal("Wrong Select Project")
	}
}
