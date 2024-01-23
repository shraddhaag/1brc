package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {

	inputFiles := find("./test_cases", ".txt")
	for _, test := range inputFiles {
		t.Run(test, func(t *testing.T) {
			output := evaluate(test + ".txt")
			assert.Equal(t, readFile(test+".out"), "{"+output+"}\n")
		})
	}
}

func readFile(input string) string {
	fileContent, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}
	return string(fileContent)
}

func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s[:len(s)-len(ext)])
		}
		return nil
	})
	return a
}
