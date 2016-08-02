package main

import (
	"fmt"
	"github.com/fheng/fh-system-dump-tool/check"
	"io"
	"path/filepath"
	"os"
	"strings"
)

func analyseTask(logFileDir string, checks []int, checkFactory check.Factory) int {
	status := 0
	for _, check := range checks {
		checker, err := checkFactory(check)
		if err != nil {
			fmt.Println(err)
			status = 1
		}

		checkFiles, err := getFilesFor(checker, logFileDir)
		if err != nil {
			fmt.Println(err)
			status = 1
		}
		res, err := checker.Execute(checkFiles)
		if err != nil {
			fmt.Println(err)
			status = 1
		}

		res.Output()
	}

	return status
}

func getFilesFor(checker check.Check, dataDir string) ([]io.Reader, error) {
	files := []string{}
	requiredFiles := checker.RequiredFiles()

	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error{
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	retFiles := []io.Reader{}

	for _, rFile := range requiredFiles {
		for _, file := range files {
			if strings.HasSuffix(file, rFile) {
				reader, err := os.Open(file)
				if err != nil {
					fmt.Println(err)
					continue
				}cd
				retFiles = append(retFiles, reader)
			}
		}
	}

	return retFiles, nil
}