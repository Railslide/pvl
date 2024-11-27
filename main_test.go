package main

import (
	"os"
	"testing"
)

type mockedFile struct {
	os.FileInfo
}

type mockedFS struct{}

func (mockedFS) Getwd() (string, error)                { return "test_directory", nil }
func (mockedFS) Stat(name string) (os.FileInfo, error) { return mockedFile{}, nil }

func TestLocateVenv(t *testing.T) {
	fs := mockedFS{}
	venv, dir, err := locateVenv(fs)

	if err != nil {
		t.Fatalf("Test encountered an error: %v", err)
	}
	if venv != ".venv" {
		t.Fatalf("Expected virtualenv name is `.venv`, but got `%s`", venv)
	}
	if dir != "test_directory" {
		t.Fatalf("Expected dir path is `test_directory`, but got `%s`", dir)
	}

}

func TestLocateVenvEnvVariable(t *testing.T) {

	fs := mockedFS{}
	os.Setenv("VIRTUAL_ENV", "test_path/test_venv")
	venv, dir, err := locateVenv(fs)

	if err != nil {
		t.Fatalf("Test encountered an error: %v", err)
	}
	if venv != "test_venv" {
		t.Fatalf("Expected virtualenv name is `test_venv`, but got `%s`", venv)
	}
	if dir != "test_path" {
		t.Fatalf("Expected dir path is `test_path`, but got `%s`", dir)
	}
}
