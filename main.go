package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func main() {
	var projectName string
	var projectType string

	setupConfig()

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
		laravelCreator(projectName)
		createVirtualServer(projectName)
		openCodeinVscode(fmt.Sprintf("%s%s", viper.GetString("laravel_path"), projectName))
		linkCreator()
		break
	case "2":
		golangCreator(projectName)
		openCodeinVscode(fmt.Sprintf("%s%s", viper.GetString("laravel_path"), projectName))
		break
	case "3":
		rustCreator(projectName)
		openCodeinVscode(fmt.Sprintf("%s%s", viper.GetString("laravel_path"), projectName))
		break
	default:
		log.Fatal("Wrong Select Project")
	}

	//make virtual server
	//
}

func setupConfig() {
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

func laravelCreator(projectName string) {
	os.Chdir(viper.GetString("laravel_path"))
	cmd := exec.Command("composer", "create-project", "laravel/laravel", projectName, "--prefer-dist")

	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

func golangCreator(projectName string) {
	os.Chdir(viper.GetString("golang_path"))
	//run mkdir projectName && go mod init projectName && create main.go >> package main func main(){}
}

func rustCreator(projectName string) {
	os.Chdir(viper.GetString("rust_path"))
	cmd := exec.Command("cargo", "new", projectName, "--bin")

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func createVirtualServer(projectName string) {
	//open sample.test.conf >> replace sample with project name >> save file to /etc/nginx/conf.d
	//open /etc/hosts >> append 127.0.0.1 projectname.test >>save
	//reload nginx
}

func openCodeinVscode(projectPath string) {
	os.Chdir(projectPath)
	cmd := exec.Command("code", ".")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func linkCreator() string {
	return "link"
}
