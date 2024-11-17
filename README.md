# pvl

pvl (pyright virtualenv locator) is a small cli utility to automatically create a pyrightconfig.json file,
so that pyright can pick up the correct virtual environment for your project.

## Installation

You can install pvl with go install
```
go install github.com/railslide/pvl@latest
```

## Virtualenv detection priority

pvl checks for virtualenv environments according to the following order

- Currently active virtualenv
- `.venv` folder in current working directory

In other words, if pvl is running from inside a virtualenv it will pick the active virtualenv
independently on whether there is a `.venv` folder or not.
