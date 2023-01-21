package projectcreator

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func GolangCreator(projectName string) {
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
