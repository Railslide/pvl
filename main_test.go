package main

import (
	"io/fs"
	"os"
	"testing"
)

type mockedFile struct {
	os.FileInfo
}

type mockedFS struct{}

func (mockedFS) Getwd() (string, error) {
	return "test_directory", nil
}

func (mockedFS) Stat(name string) (os.FileInfo, error) {
	return mockedFile{}, nil
}

func (mockedFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
  return os.ErrNotExist
}


func TestExtractVent(t *testing.T) {
	fs := mockedFS{}

	tests := []struct {
		name        string
		venv        string
		envVarName  string
		envVarValue string
		wantVenv    string
		wantDir     string
	}{
		{
			name:     "Full path",
			venv:     "/hello/world/test_venv",
			wantVenv: "test_venv",
			wantDir:  "/hello/world",
		},
		{
			name:     "Relative path",
			venv:     "test_vent",
			wantVenv: "test_vent",
			wantDir:  "test_directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			venv, dir, err := extractVenv(fs, tt.venv)

			if err != nil {
				t.Errorf("got %v, expected nil", err)
			}

			if venv != tt.wantVenv {
				t.Errorf("got: %s, want: %s", venv, tt.wantVenv)
			}

			if dir != tt.wantDir {
				t.Errorf("got: %s, want: %s", dir, tt.wantDir)
			}
		})
	}
}

func TestLocateVenv(t *testing.T) {
	fs := mockedFS{}

	tests := []struct {
		name        string
		targetDir   string
		envVarName  string
		envVarValue string
		wantVenv    string
		wantDir     string
	}{
		{
			name:        "VIRTUAL_ENV env variable is set",
			envVarName:  "VIRTUAL_ENV",
			envVarValue: "/test_path/test_venv",
			wantVenv:    "test_venv",
			wantDir:     "/test_path",
		},
		{
			name:     ".venv folder in current working directory",
			wantVenv: ".venv",
			wantDir:  "test_directory",
		},
		{
			name:      "user-specified folder",
			targetDir: "test_directory",
			wantVenv:  "test_directory",
			wantDir:   "test_directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVarName != "" {
				os.Setenv(tt.envVarName, tt.envVarValue)

				t.Cleanup(func() { os.Unsetenv(tt.envVarName) })
			}

			venv, dir, err := locateVenv(fs, tt.targetDir)

			if err != nil {
				t.Errorf("got %v, expected nil", err)
			}

			if venv != tt.wantVenv {
				t.Errorf("got: %s, want: %s", venv, tt.wantVenv)
			}

			if dir != tt.wantDir {
				t.Errorf("got: %s, want: %s", dir, tt.wantDir)
			}

		})
	}
}

func TestCreateConfigFileAlreadyExisting(t *testing.T) {
	fs := mockedFS{}

	err := createConfigFile(fs, "test_venv", "test_path")

	if err == nil || err.Error() != "Config file already exists" {
		t.Errorf("Test did not throw an error for already existing config file")
	}
}
