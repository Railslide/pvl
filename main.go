package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/railslide/pvl/internal"
)


func main() {
	log.SetFlags(0)
	log.SetPrefix("pvl: ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of pvl:\n")
		fmt.Fprintf(os.Stderr, "\tpvl\n")
		fmt.Fprintf(os.Stderr, "\tpvl --path [PATH]\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	userPath := flag.String("path", "", "Custom path to the virtualenv (optional)")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "error: unrecognized command\n\n")
		flag.Usage()
		os.Exit(1)
	}

	internal.CreateConfig(*userPath)
}
