package main

import (
	"encoding/json"
	"fmt"
	"github.com/fheng/fh-system-dump-tool/check"
	"os"
	"path/filepath"
	"strings"
	"io"
)


// outputResults will take the results returned from analyseTask and output them
// honouring the format and pretty strings. Current valid format is only JSON
func outputResults(results map[string][]*check.Result, format string, pretty bool) error {
	var output []byte
	var err error
	if pretty {
		output, err = json.MarshalIndent(results, "", "    ")
	} else {
		output, err = json.Marshal(results)
	}

	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}

// The analyse task takes an array of files names, an array of check IDs and a
// factory to build checks from. It will execute the checks on any relevant
// files in the list of files and return a map of results or else an error
func analyseTask(files []string, checks []check.CheckTask, checkFactory check.Factory) (map[string][]*check.Result, error) {
	fileFactory := getFilesFactory(files)

	results := map[string][]*check.Result{"results": {}}

	for _, check := range checks {
			res = check(fileFactory)
			if err != nil {
				return nil, err
			}
		}
		results["results"] = append(results["results"], checker.GetResult())
	}

	return results, nil
}

// getFilesFor takes a check with a RequiredFiles method and an array of file
// names and returns the paths of any files which match the checks
// RequiredFiles.
func getFilesFactory(files []string) check.CheckFileFactory {
	return func(requiredFiles []string) []io.Reader {
		retFiles := []io.Reader{}
		for _, rFile := range requiredFiles {
			for _, file := range files {
				if strings.HasSuffix(file, rFile) {
					retFiles = append(retFiles, file)
				}
			}
		}
		return retFiles
	}
}

// This function recursively lists the relative path of all files in the
// provided directory and all sub-directories, returned as an array of
// strings
func listAllFiles(dir string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
