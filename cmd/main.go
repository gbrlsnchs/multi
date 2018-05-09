package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/gbrlsnchs/multi"
)

func main() {
	var (
		f       []byte
		err     error
		cfgFile string
		noColor bool
	)

	flag.StringVar(&cfgFile, "config", "multi.json", "path for a configuration file, defaults to \"multi.json\"")
	flag.BoolVar(&noColor, "no-color", false, "disable colorful output")
	flag.Parse()

	if f, err = ioutil.ReadFile(cfgFile); err != nil {
		log.Fatal(err)
	}

	log.Fatal(multi.Start(f, noColor))
}
