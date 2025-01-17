package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
)

type fileSystem interface {
	Getwd() (string, error)
	Stat(name string) (os.FileInfo, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type osFS struct{}

func (osFS) Getwd() (string, error)                { return os.Getwd() }
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
func (osFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

type PyrightConfig struct {
	VenvName string `json:"venv"`
	VenvPath string `json:"venvPath"`
}

func locateVenv(fs fileSystem) (venvName, venvPath string, err error) {
	envVar, ok := os.LookupEnv("VIRTUAL_ENV")
	if ok {
		venvPath, venvName := path.Split(envVar)
		return venvName, path.Clean(venvPath), nil
	}

	cwd, err := fs.Getwd()
	if err != nil {
		return "", "", errors.New("Cannot get path of current working directory")
	}

	venvName = ".venv"
	localVenvDir := path.Join(cwd, venvName)
	if _, err := fs.Stat(localVenvDir); err != nil {
		return "", "", errors.New("Cannot find a virtualenv for the project")
	}

	return venvName, cwd, nil
}

func createConfigFile(fs fileSystem, venvName, venvPath string) error {
	filename := "pyrightconfig.json"
	config := PyrightConfig{
		VenvName: venvName,
		VenvPath: venvPath,
	}

	fileContent, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return errors.New("Error while creating file content")
	}

	if _, err := fs.Stat(filename); !errors.Is(err, os.ErrNotExist) {
		return errors.New("Config file already exists")
	}

        fileContent = append(fileContent, '\n')

	err = fs.WriteFile(filename, fileContent, 0644)
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

	fs := osFS{}
	venvName, venvPath, err := locateVenv(fs)
	if err != nil {
		log.Fatal(err)
	}

	err = createConfigFile(fs, venvName, venvPath)
	if err != nil {
		log.Fatal(err)
	}
}
