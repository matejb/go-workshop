package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

func main() {

}

type cliParams struct {
	path string
	out  string
}

// parseFlags will parse cli program arguments into internal structure for later use
func parseFlags() (params cliParams) {

	// needed info: https://golang.org/pkg/flag/#StringVar

	flag.StringVar(&params.path, "list", "", "Path to dir containg css files")
	flag.StringVar(&params.out, "out", "", "Filename of destination css file")

	// needed info: https://golang.org/pkg/flag/#Parse

	flag.Parse()

	return
}

// List will read list of css files from list json file
func List(listFile string) (cssFilePaths []string, err error) {

	// needed info: https://golang.org/pkg/io/ioutil/#ReadFile

	content, err := ioutil.ReadFile(listFile)
	if err != nil {
		return cssFilePaths, err
	}

	// needed info: https://golang.org/pkg/encoding/json/#Unmarshal

	err = json.Unmarshal(content, &cssFilePaths)
	if err != nil {
		return cssFilePaths, err
	}

	return
}
