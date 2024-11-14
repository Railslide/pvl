package main

import (
	"log"
	"os"
	"path"
)

func locateVenv() string {
  envVar, hasValue := os.LookupEnv("VIRTUAL_ENV")
  if hasValue {
    return envVar
  }

  cwd, err := os.Getwd()
  if err != nil {
    log.Fatal("Cannot get path of current working directory")
  }

  localVenvDir := path.Join(cwd, ".venv")
  if _, err := os.Stat(localVenvDir); err != nil {
    log.Fatal("Cannot find a virtualenv for the project")
  }

  return localVenvDir
}

func createConfigFile(venvPath string) {

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

  venvPath := locateVenv()
  println(venvPath)

}

