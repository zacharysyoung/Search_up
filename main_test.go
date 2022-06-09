package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

const tempDirName = "Search-up_temp"
const tempFileName = "temp-file.txt"

func TestNonExistingFile(t *testing.T) {
	rootTempDir, tempDirPaths := setupTempFS()
	defer cleanTempFS(rootTempDir)

	deepestPath := tempDirPaths[0]

	foundPaths, err := searchUp(deepestPath, "BOGUS-FILE", false)
	if err != nil {
		t.Error(err)
	}

	if len(foundPaths) != 0 {
		t.Errorf("Expected to find 0 occurrences of \"bogus-file\", found %q\n", foundPaths)
	}
}

func TestFindSingleDir(t *testing.T) {
	rootTempDir, tempDirPaths := setupTempFS()
	defer cleanTempFS(rootTempDir)

	// "Work up" the directory tree from the bottom to the top; should see less-and-less results
	tempCount := len(tempDirPaths)
	for i := 0; i < tempCount; i++ {
		tempPaths := tempDirPaths[i:tempCount]
		deepestPath := tempPaths[0]

		foundPaths, err := searchUp(deepestPath, tempFileName, true)
		if err != nil {
			t.Error(err)
		}

		if len(foundPaths) != 1 {
			t.Errorf("Expected to find 1 occurrence of temp-file in dir path \"%s\", instead found %q\n", deepestPath, foundPaths)
		}

		if foundPaths[0] != deepestPath {
			t.Errorf("Expected to find \"%s\", instead found \"%s\"\n", deepestPath, foundPaths[0])
		}
	}

}

func TestFindMultipleDirs(t *testing.T) {
	rootTempDir, tempDirPaths := setupTempFS()
	defer cleanTempFS(rootTempDir)

	// "Work up" the directory tree from the bottom to the top; should see less-and-less results
	tempCount := len(tempDirPaths)
	for i := 0; i < tempCount; i++ {
		tempPaths := tempDirPaths[i:tempCount]
		deepestPath := tempPaths[0]

		foundPaths, err := searchUp(deepestPath, tempFileName, false)
		if err != nil {
			t.Error(err)
		}

		if len(foundPaths) != len(tempPaths) {
			t.Errorf("Expected to find 3 occurrence of temp-file in dir paths %q, instead found %q\n", tempDirPaths, foundPaths)
		}

		for j := 0; j < len(tempPaths); j++ {
			want := tempPaths[j]
			got := foundPaths[j]
			if got != want {
				t.Errorf("At sub-dir level %d, expected to find \"%s\", instead found \"%s\"\n", j, want, got)
			}
		}
	}
}

func TestBadPath(t *testing.T) {
	_, err := searchUp("a/bad/path", "irrelevant-name", false)
	if err == nil {
		t.Error("Wanted an error for searchUp(\"a/bad/path\", ...), but got <nil>")
	}
}

// setupTempFS creates a temp directory structure, three sub-dirs deep: rootTempDir/SubA/SubB/SubC.
// Each sub-dir has a copy of tempFileName, e.g.: rootTempDir/SubA/tmp-file.txt, etc...
// Returns a list of each temp-file in its sub-dir.
func setupTempFS() (string, []string) {

	rootTempDirPath, err := os.MkdirTemp("", tempDirName)
	if err != nil {
		log.Fatal("could not make rootTempDir", tempDirName, err)
	}
	subDirPath := rootTempDirPath

	subDirNames := []string{"SubA", "SubB", "SubC"}
	subDirPaths := make([]string, len(subDirNames))

	for i, subDirName := range subDirNames {
		subDirPath = filepath.Join(subDirPath, subDirName)

		err := os.Mkdir(subDirPath, 0750)
		if err != nil {
			log.Fatal("could not create subDir", subDirPath, err)
		}

		subFilePath := filepath.Join(subDirPath, tempFileName)
		f, err := os.Create(subFilePath)
		if err != nil {
			log.Fatal("could not create subFile", subFilePath, err)
		}
		f.Close()

		subDirPaths[len(subDirNames)-1-i] = subDirPath
	}
	return rootTempDirPath, subDirPaths
}

func cleanTempFS(rootTempDir string) {
	os.RemoveAll(rootTempDir)
}
