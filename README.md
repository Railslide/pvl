# pvl

pvl (pyright virtualenv locator) is a small cli utility to automatically create a pyrightconfig.json file,
so that pyright can pick up the correct virtual environment for your project.

## Installation

You can install pvl with go install
```
go install github.com/railslide/pvl@latest
```

## Usage
```
$ pvl --help

Usage of pvl:
	pvl
	pvl --path [PATH]
Flags:
  -path string
       Custom path to the virtualenv (optional)
```

## Virtualenv detection priority

pvl checks for virtualenv environments according to the following order:

- Path passed as a command line flag
- Currently active virtualenv
- `.venv` folder in current working directory
