package main

import (
	"log"

	"github.com/railslide/pvl/cmd"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("pvl: ")

	cmd.Execute()
}
