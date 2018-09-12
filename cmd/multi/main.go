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
	noColor := flag.Bool("nocolor", false, "disable colorful output")
	noPid := flag.Bool("nopid", false, "disable printing PID along with the output")
	noDate := flag.Bool("nodate", false, "disable logging date")
	noTime := flag.Bool("notime", false, "disable logging time")
	silent := flag.Bool("silent", false, "disable output")
	stdout := flag.String("stdout", "", "log stdout to a file")
	stderr := flag.String("stderr", "", "log stderr to a file")
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

	if !*silent {
		tl.Stderr = []io.Writer{os.Stderr}
		tl.Stdout = []io.Writer{os.Stdout}
	}
	var f *os.File
	if *stderr != "" {
		if f, err = os.OpenFile(*stderr, os.O_APPEND, os.ModeAppend); err != nil {
			if !os.IsNotExist(err) {
				log.Fatal(err)
			}
			if f, err = os.Create(*stderr); err != nil {
				log.Fatal(err)
			}
		}
		tl.Stderr = append(tl.Stderr, f)
	}
	if *stdout != "" {
		if f, err = os.OpenFile(*stdout, os.O_APPEND, os.ModeAppend); err != nil {
			if !os.IsNotExist(err) {
				log.Fatal(err)
			}
			if f, err = os.Create(*stdout); err != nil {
				log.Fatal(err)
			}
		}
		tl.Stdout = append(tl.Stdout, f)
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
