package main

import (
	"encoding/json"
	"io/ioutil"
)

func main() {

}

// Flag will read cli flag (parametar) and return it's value
func Flag(flagName string) (value string) {
	return
}

// List will read list of css files from list json file
func List(listFile string) (cssFilePaths []string, err error) {

	// hint: https://golang.org/pkg/io/ioutil/#ReadFile

	content, err := ioutil.ReadFile(listFile)
	if err != nil {
		return cssFilePaths, err
	}

	// hint: https://golang.org/pkg/encoding/json/#Unmarshal

	err = json.Unmarshal(content, &cssFilePaths)
	if err != nil {
		return make([]string, 0), err
	}

	return
}
