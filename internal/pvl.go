package internal

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

func locateVenv(fs fileSystem, userPath string) (venvName, venvPath string, err error) {
	var targetPath string

	if userPath == "" {
		envVar, ok := os.LookupEnv("VIRTUAL_ENV")
		if !ok {
			targetPath = ".venv"
		} else {
			targetPath = envVar
		}

	} else {
		targetPath = userPath
	}

	venvName, venvPath, err = extractVenv(fs, targetPath)
	if err != nil {
		return "", "", err
	}

	venvFullPath := path.Join(venvPath, venvName)
	if _, err := fs.Stat(venvFullPath); err != nil {
		return "", "", fmt.Errorf("cannot find a virtualenv at %s", venvFullPath)
	}

	return venvName, venvPath, nil
}

func writeConfigFile(fs fileSystem, venvName, venvPath, filename string) error {
	config := PyrightConfig{
		VenvName: venvName,
		VenvPath: venvPath,
	}

	fileContent, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return errors.New("error while creating file content")
	}

	fileContent = append(fileContent, '\n')

	err = fs.WriteFile(filename, fileContent, 0644)
	if err != nil {
		return errors.New("error while writing file")
	}
	return nil
}

func CreateConfig(userPath string) {
	fs := osFS{}
	destination := "pyrightconfig.json"

	if _, err := fs.Stat(destination); !errors.Is(err, os.ErrNotExist) {
		log.Fatal("config file already exists")
	}

	venvName, venvPath, err := locateVenv(fs, userPath)
	if err != nil {
		log.Fatal(err)
	}

	err = writeConfigFile(fs, venvName, venvPath, destination)
	if err != nil {
		log.Fatal(err)
	}
}
