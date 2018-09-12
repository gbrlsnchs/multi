package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gbrlsnchs/multi"
)

func main() {
	cfg := flag.String("config", "multi.toml", "path for a configuration file, defaults to multi.toml")
	noColor := flag.Bool("nocolor", false, "disable colorful output")
	flag.Parse()

	b, err := ioutil.ReadFile(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	var tl multi.TaskList
	ext := filepath.Ext(*cfg)
	switch ext {
	case ".json":
		err = json.Unmarshal(b, &tl)
	case ".toml":
		err = toml.Unmarshal(b, &tl)
	default:
		err = fmt.Errorf("multi: unsupported config extension %s", ext)
	}
	if err != nil {
		log.Fatal(err)
	}
	tl.Stderr = []io.Writer{os.Stderr}
	tl.Stdout = []io.Writer{os.Stdout}
	if err = tl.Start(!*noColor); err != nil {
		log.Fatal(err)
	}
}
