package virtualHost

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Create(projectPath string, projectName string) {
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
	var content string
	for scanner.Scan() {
		addressReplace := strings.Replace(scanner.Text(), "address", projectPath, -1)
		projectNameReplace := strings.Replace(addressReplace, "sample", projectName, -1)
		content = fmt.Sprintf("%s%s", projectNameReplace, "\n")
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

	cmd := exec.Command("sudo", "mv", sourceFile, destFile)
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
	content = fmt.Sprintf("%s%s%s", "127.0.0.1\t", projectName, ".test\n")
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
