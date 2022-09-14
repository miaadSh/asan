package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

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
		linkCreator(projectName)
		break
	case "2":
		golangCreator(projectName)
		openCodeinVscode(fmt.Sprintf("%s%s", viper.GetString("golang_path"), projectName))
		break
	case "3":
		rustCreator(projectName)
		openCodeinVscode(fmt.Sprintf("%s%s", viper.GetString("rust_path"), projectName))
		break
	default:
		log.Fatal("Wrong Select Project")
	}
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
	currentPath, _ := os.Getwd()
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

	dir := fmt.Sprintf("%s%s", viper.GetString("laraavel_path"), projectName)
	os.Chdir(dir)

	cmd = exec.Command("sudo", "chgrp", "-R", "www-data", "storage")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("OUT:", string(out))
	}

	cmd = exec.Command("sudo", "chgrp", "-R", "www-data", "bootstrap/cache")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("OUT:", string(out))
	}

	cmd = exec.Command("sudo", "chmod", "-R", "775", "storage")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("OUT:", string(out))
	}

	cmd = exec.Command("sudo", "chmod", "-R", "775", "bootstrap/cache")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("OUT:", string(out))
	}

	os.Chdir(currentPath)
}

func golangCreator(projectName string) {
	path := fmt.Sprintf("%s%s", viper.GetString("golang_path"), projectName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(projectName, " Directory created")
	} else {
		fmt.Println("Directory already exists")
	}

	source, sourceError := os.Open("./stub/main.go.stub")
	if sourceError != nil {
		log.Fatal(sourceError)
	}
	defer source.Close()

	os.Chdir(fmt.Sprintf("%s%s", viper.GetString("golang_path"), projectName))

	target, targetError := os.OpenFile("main.go", os.O_RDWR|os.O_CREATE, 0755)
	if targetError != nil {
		log.Fatal(targetError)
	}
	defer target.Close()

	_, copyError := io.Copy(target, source)
	if copyError != nil {
		log.Fatal(copyError)
	}

	name := fmt.Sprintf("%s%s", "github.com/", projectName)
	cmd := exec.Command("go", "mod", "init", name)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("create golang projects")
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
	source, sourceError := os.OpenFile("./stub/sample.test.conf.stub", os.O_RDWR|os.O_CREATE, 0755)
	if sourceError != nil {
		log.Fatal(sourceError)
	}
	defer source.Close()

	newFile, createError := os.Create("./stub/move.conf")
	if createError != nil {
		log.Fatal(createError)
	}
	defer newFile.Close()

	scanner := bufio.NewScanner(source)
	// buf := make([]byte, 0, 1024)
	// scanner.Buffer(buf, 256*1024)
	for scanner.Scan() {
		content := fmt.Sprintf("%s%s", strings.Replace(scanner.Text(), "sample", projectName, -1), "\n")
		_, writeError := newFile.WriteString(content)
		if writeError != nil {
			log.Fatal(writeError, 123)
		}
	}

	error1 := scanner.Err()
	if error1 != nil {
		log.Fatal(error1, 345)
	}

	sourceFile := "./stub/move.conf"
	destFile := fmt.Sprintf("%s%s%s", "/etc/nginx/conf.d/", projectName, ".test.conf")

	cmd := exec.Command("sudo", "cp", "-p", sourceFile, destFile)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("OUT:", string(out))
	}

	cmd = exec.Command("sudo", "cat", "/etc/hosts")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("OUT:", string(out))
	}

	hostFile, hostCreateError := os.Create("./stub/hosts")
	if hostCreateError != nil {
		log.Fatal(hostCreateError)
	}
	defer hostFile.Close()

	_, writeError := hostFile.WriteString(string(out))
	if writeError != nil {
		log.Fatal(writeError)
	}
	content := fmt.Sprintf("%s%s%s", "127.0.0.1\t", projectName, ".test\n")
	_, writeError = hostFile.WriteString(content)
	if writeError != nil {
		log.Fatal(writeError)
	}

	error1 = scanner.Err()
	if error1 != nil {
		log.Fatal(error1, 345)
	}
	sourceFile = "./stub/hosts"
	destFile = "/etc/hosts"

	cmd = exec.Command("sudo", "mv", sourceFile, destFile)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("OUT:", string(out))
	}

	cmd = exec.Command("sudo", "systemctl", "reload", "nginx.service")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("OUT:", string(out))
	}
}

func openCodeinVscode(projectPath string) {
	os.Chdir(projectPath)
	cmd := exec.Command("code", ".")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func linkCreator(projectName string) {
	url := fmt.Sprintf("%s%s%s", "http://", projectName, ".test")
	openBrowser(url)
	fmt.Println(url)
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
