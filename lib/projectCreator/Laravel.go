package projectcreator

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func LaravelCreator(projectName string) {
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

	dir := fmt.Sprintf("%s%s", viper.GetString("laravel_path"), projectName)
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
