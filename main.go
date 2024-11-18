package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
)

type PyrightConfig struct {
	VenvName string `json:"venv"`
	VenvPath string `json:"venvPath"`
}

func locateVenv() (venvName, venvPath string, err error) {
	envVar, hasValue := os.LookupEnv("VIRTUAL_ENV")
	if hasValue {
		venvPath, venvName := path.Split(envVar)
		return venvName, venvPath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", "", errors.New("Cannot get path of current working directory")
	}

	venvName = ".venv"
	localVenvDir := path.Join(cwd, venvName)
	if _, err := os.Stat(localVenvDir); err != nil {
		return "", "", errors.New("Cannot find a virtualenv for the project")
	}

	return venvName, cwd, nil
}

func createConfigFile(venvName, venvPath string) error {
	filename := "pyrightconfig.json"
	pyrightConfig := PyrightConfig{
		VenvName: venvName,
		VenvPath: venvPath,
	}

	fileContent, err := json.MarshalIndent(pyrightConfig, "", "    ")
	if err != nil {
		return errors.New("Error while creating file content")
	}

	if _, err := os.Stat(filename); err == nil {
		return errors.New("Config file already exists")
	}

	err = os.WriteFile("pyrightconfig.json", fileContent, 0644)
	if err != nil {
		return errors.New("Error while writing file")
	}
	return nil
}

func main() {

	if len(os.Args) > 1 {
		if os.Args[1] == "--help" {
			println(
				"Pvl, the Pyright Virtualenv Locator\n\n" +
					"It locates the virtualenv and adds the " +
					"path to it to pyrightconfig.json",
			)
			os.Exit(0)
		} else {
			println("Unrecognized command")
			os.Exit(1)
		}
	}

	venvName, venvPath, err := locateVenv()
	if err != nil {
		log.Fatal(err)
	}

	err = createConfigFile(venvName, venvPath)
	if err != nil {
		log.Fatal(err)
	}
}
