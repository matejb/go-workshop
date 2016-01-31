package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// main is the function that will be called when program starts, when main function exists program exists
func main() {

	params := parseFlags()

	cssPaths, err := list(params.list)
	if err != nil {
		log.Fatal(err)
	}

	err = merge(cssPaths, params.out)
	if err != nil {
		log.Fatal(err)
	}

	if params.watch {
		err := watch(cssPaths, params.out)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type cliParams struct {
	list  string
	out   string
	watch bool
}

// parseFlags will parse cli program arguments into internal structure for later use
func parseFlags() (params cliParams) {

	// needed info: https://golang.org/pkg/flag/#StringVar

	flag.StringVar(&params.list, "list", "", "Path to dir containg css files")
	flag.StringVar(&params.out, "out", "", "Filename of destination css file")
	flag.BoolVar(&params.watch, "watch", false, "Enables watch mode that automatically rebuilds destination css file if any of source css files changes")

	// needed info: https://golang.org/pkg/flag/#Parse

	flag.Parse()

	return
}

// list will read list of css files from list json file
func list(listFile string) (cssFilePaths []string, err error) {

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

// merge will merge css files into one big new file, if merged file exists it will be overwritten
func merge(cssFilePaths []string, mergedFile string) (err error) {

	// needed info:
	// https://golang.org/pkg/os/#Create
	// https://golang.org/pkg/os/#File.Close
	// https://golang.org/pkg/os/#Open
	// https://golang.org/pkg/io/#Copy

	out, err := os.Create(mergedFile)
	if err != nil {
		return
	}
	defer out.Close()

	for _, path := range cssFilePaths {
		in, err := os.Open(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, in)
		if err != nil {
			in.Close()
			return err
		}

		in.Close()
	}

	return nil
}

// watch will watch changes in cssFilePaths files and rebuild mergedFile
func watch(cssFilePaths []string, mergedFile string) (err error) {

	fmt.Println("Watch mode: started")

	rebuildCh := make(chan bool)
	errCh := make(chan error)
	cleanupCh := make(chan struct{})

	watchSingleFile := func(path string, rebuildCh chan bool, errCh chan error, cleanupCh chan struct{}) {
		initStat, err := os.Stat(path)
		if err != nil {
			select {
			case errCh <- err:
				return
			case <-cleanupCh:
				return
			}
		}

		fmt.Printf("Watch mode: added %q to watch list\n", path)

		for {
			stat, err := os.Stat(path)
			if err != nil {
				select {
				case errCh <- err:
					return
				case <-cleanupCh:
					return
				}
			}

			if stat.Size() != initStat.Size() || stat.ModTime() != initStat.ModTime() {
				select {
				case rebuildCh <- true:
					initStat = stat
				case <-cleanupCh:
					return
				}
			}

			time.Sleep(50 * time.Millisecond)
		}
	}

	for _, path := range cssFilePaths {
		go watchSingleFile(path, rebuildCh, errCh, cleanupCh)
	}

	for {
		select {

		case <-rebuildCh:
			err := merge(cssFilePaths, mergedFile)
			if err != nil {
				close(cleanupCh)
				return err
			}
			fmt.Printf("Watch mode: rebuilded %q\n", mergedFile)

		case err := <-errCh:
			close(cleanupCh)
			return err
		}
	}

}
