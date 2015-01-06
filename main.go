package main

import (
	"flag"
	"log"
	"os"
	"runtime"
)

var show_h bool
var show_help bool
var show_version bool

func init() {
	flag.BoolVar(&show_h, "h", false, "show help message and exit(0)")
	flag.BoolVar(&show_help, "help", false, "show help message and exit(0)")
	flag.BoolVar(&show_version, "version", false, "show version info and exit(0)")
}

func main() {
	log.Println("paral")
	log.Printf("- version %v (build %d)\n", Version(), BuildNumber())
	log.Printf("- Go Version: %v %v/%v\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	log.Printf("- NumCPU(): %v\n", runtime.NumCPU())
	log.Printf("- GOMAXPROCS(): %v\n", runtime.GOMAXPROCS(0))

	// command line flags:
	flag.Parse()

	if show_version {
		os.Exit(0)
	}

	if show_h || show_help {
		flag.Usage()
		os.Exit(0)
	}

}
