package projectcreator

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func RustCreator(projectName string) {
	os.Chdir(viper.GetString("rust_path"))
	cmd := exec.Command("cargo", "new", projectName, "--bin")

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}
