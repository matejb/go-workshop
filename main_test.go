package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

func TestFlags(t *testing.T) {

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	_ = parseFlags()

	if !flag.Parsed() {
		t.Error("Expected cli flags to be parsed!")
	}

	if flag.Lookup("list") == nil {
		t.Errorf("Expected cli flag %q to be readed", "list")
	}

	if flag.Lookup("out") == nil {
		t.Errorf("Expected cli flag %q to be readed", "out")
	}
}

func TestList(t *testing.T) {
	expected := []string{
		"test_resources/first.css",
		"test_resources/second.css",
		"test_resources/third.css",
	}

	listPath := "test_resources" + string(os.PathSeparator) + "list.js"

	result, err := list(listPath)
	if err != nil {
		t.Fatalf("List loading failed with error %s", err)
	}

	if len(result) != len(expected) {
		t.Fatal("Paths are not readed from json file.")
	}

	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Errorf("For path %d, expeced %s but recoved %s", i, expected[i], result[i])
		}
	}
}

func TestMarge(t *testing.T) {
	paths := []string{
		"test_resources" + string(os.PathSeparator) + "first.css",
		"test_resources" + string(os.PathSeparator) + "second.css",
		"test_resources" + string(os.PathSeparator) + "third.css",
	}

	tempFile, err := ioutil.TempFile("test_resources", "cssMergeTest")
	if err != nil {
		t.Fatal(err)
	}

	stat, err := tempFile.Stat()
	if err != nil {
		tempFile.Close()
		t.Fatal(err)
	}

	out := "test_resources" + string(os.PathSeparator) + stat.Name()

	tempFile.Close()

	defer func() {
		os.Remove(out)
	}()

	err = merge(paths, out)
	if err != nil {
		t.Errorf("Merge failed with error %q", err)
	}

	outContent, err := ioutil.ReadFile(out)
	if err != nil {
		t.Errorf("Reading of out file failed with error %q", err)
	}

	var expected []byte
	for _, path := range paths {
		cssContent, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("Test setup failed while processing path %q with error %q", path, err)
		}
		expected = append(expected, cssContent...)
	}

	if string(outContent) != string(expected) {
		t.Errorf("Expected %q recived %q", expected, outContent)
	}
}

// func TestFaze1(t *testing.T) {
// 	outPath := "test_resources" + string(os.PathSeparator) + "merged.css"
// 	inPath := "test_resources" + string(os.PathSeparator) + "list.js"

// 	defer func() {
// 		if _, err := os.Stat(outPath); err == nil {
// 			err := os.Remove(outPath)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 		}
// 	}()

// 	cmd := exec.Command("go", "run", "main.go", "--list", inPath, "--out", outPath)
// 	err := cmd.Run()
// 	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
// 		t.Fatalf("Faze 1 test process ran with err %v, want exit status 0", err)
// 	}

// 	expected := `h1 {color: red;}
// p {color: purple;}
// strong {color: yellow;}`

// 	result, err := ioutil.ReadFile(outPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(result) != expected {
// 		t.Fatalf("Faze 1 fails, got %s", result)
// 	}
// }
