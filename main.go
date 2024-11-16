package main

import (
        "encoding/json"
	"log"
	"os"
	"path"
)

type PyrightConfig struct {
  VenvName string
  VenvPath string
}

func locateVenv() (venvName, venvPath string) {
  envVar, hasValue := os.LookupEnv("VIRTUAL_ENV")
  if hasValue {
    venvPath, venvName := path.Split(envVar)
    return venvName, venvPath
  }

  cwd, err := os.Getwd()
  if err != nil {
    log.Fatal("Cannot get path of current working directory")
  }

  venvName = ".venv"
  localVenvDir := path.Join(cwd, venvName)
  if _, err := os.Stat(localVenvDir); err != nil {
    log.Fatal("Cannot find a virtualenv for the project")
  }

  return venvName, cwd
}

func createConfigFile(venvName, venvPath string) {
  jsonBody := PyrightConfig{
    VenvName: venvName,
    VenvPath: venvPath,
  }

  fileContent, err := json.MarshalIndent(jsonBody, "", "    ")
  if err != nil {
    log.Fatal("CHANGE ME")
  }
  println(string(fileContent))
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

  venvName, venvPath := locateVenv()
  println(venvPath, venvName)

  createConfigFile(venvName, venvPath)

}

