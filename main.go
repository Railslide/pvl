package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
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

func extractVenv(fs fileSystem, venv string) (venvName, venvPath string, err error) {
	if strings.HasPrefix(venv, "/") {
		venvPath, venv := path.Split(venv)
		return venv, path.Clean(venvPath), nil
	}

	cwd, err := fs.Getwd()
	if err != nil {
		return "", "", err
	}

	return venv, cwd, nil
}

func locateVenv(fs fileSystem, venvDir string) (venvName, venvPath string, err error) {
	var target string

	if venvDir == "" {
		envVar, ok := os.LookupEnv("VIRTUAL_ENV")
		if !ok {
			target = ".venv"
		} else {
			target = envVar
		}

	} else {
		target = venvDir
	}

	venvName, venvPath, err = extractVenv(fs, target)
	if err != nil {
		return "", "", err
	}

	venvFullPath := path.Join(venvPath, venvName)
	if _, err := fs.Stat(venvFullPath); err != nil {
		return "", "", fmt.Errorf("Cannot find a virtualenv at %s", venvFullPath)
	}

	return venvName, venvPath, nil
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
			println(`Pvl, the Pyright Virtualenv Locator\n\n"
				It locates the virtualenv and adds the
				path to it to pyrightconfig.json`,
			)
			os.Exit(0)
		} else {
			println("Unrecognized command")
			os.Exit(1)
		}
	}

	fs := osFS{}
	venvDir := "foo"
	venvName, venvPath, err := locateVenv(fs, venvDir)
	if err != nil {
		log.Fatal(err)
	}

	err = createConfigFile(fs, venvName, venvPath)
	if err != nil {
		log.Fatal(err)
	}
}
