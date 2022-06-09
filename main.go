package main

// Starting at the current directory, search up the directory tree for a name file.
//
// By default print all found instances of file, optionally stop on first occurence.
//
// / (root) is not searched.
//
// Inspired by https://stackoverflow.com/a/19011599/246801

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var onlyFirst bool
var startPath string
var fileName string

func initCmd() {
	flag.BoolVar(&onlyFirst, "f", false, "Print the first occurrence of filename and stop walking up the tree.")

	flag.Parse()

	tailArgs := flag.Args()
	if len(tailArgs) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [-f] FILENAME\n", os.Args[0])
		os.Exit(1)
	}
	fileName = tailArgs[0]

	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not get working directory: ", err)
		os.Exit(1)
	}
	startPath = path
}

func main() {

	initCmd()

	foundPaths, err := searchUp(startPath, fileName, onlyFirst)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	for _, foundPath := range foundPaths {
		fmt.Println(foundPath)
	}
}

// Does not search root (/)
func searchUp(path string, fname string, onlyFirst bool) ([]string, error) {
	var foundPaths []string
	// foo
	for {
		found, err := searchDirForFile(path, fname)
		if err != nil {
			return nil, fmt.Errorf("could not read DirectoryEntries (list) \"%s\": %q", path, err)
		}

		if found {
			foundPaths = append(foundPaths, path)
			if onlyFirst {
				break
			}
		}

		path = filepath.Dir(path)

		if path == "/" {
			break
		}
	}

	return foundPaths, nil
}

func searchDirForFile(dir string, file string) (bool, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.Name() == file {
			return true, nil
		}
	}

	return false, nil
}
