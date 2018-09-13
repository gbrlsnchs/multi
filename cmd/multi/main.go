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
	cfg := flag.String("config", "multi.toml", "choose another config file")
	debug := flag.String("debug", "", "log errors to a file")
	noColor := flag.Bool("nocolor", false, "disable colorful output")
	noPid := flag.Bool("nopid", false, "disable printing PID along with the output")
	noDate := flag.Bool("nodate", false, "disable logging date")
	noTime := flag.Bool("notime", false, "disable logging time")
	silent := flag.Bool("silent", false, "disable output")
	stdout := flag.String("stdout", "", "log stdout to a file")
	stderr := flag.String("stderr", "", "log stderr to a file")
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("multi: ")
	if *debug != "" {
		log.SetOutput(io.MultiWriter(os.Stderr, openOrCreate(*debug)))
	}

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

	if !*silent {
		tl.Stderr = []io.Writer{os.Stderr}
		tl.Stdout = []io.Writer{os.Stdout}
	}
	if *stderr != "" {
		tl.Stderr = append(tl.Stderr, openOrCreate(*stderr))
	}
	if *stdout != "" {
		tl.Stdout = append(tl.Stdout, openOrCreate(*stdout))
	}
	tl.Flags = log.Ldate | log.Ltime | multi.Mcolor | multi.Mpid
	if *noColor {
		tl.Flags ^= multi.Mcolor
	}
	if *noPid {
		tl.Flags ^= multi.Mpid
	}
	if *noDate {
		tl.Flags ^= log.Ldate
	}
	if *noTime {
		tl.Flags ^= log.Ltime
	}

	if err = tl.Start(); err != nil {
		log.Fatal(err)
	}
}

func openOrCreate(fpath string) *os.File {
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
