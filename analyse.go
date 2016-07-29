package main

import (
	"fmt"
	"github.com/fheng/fh-system-dump-tool/check"
)

func analyseTask(logFileDir string, checks []int, checkFactory check.Factory) int {
	status := 0
	for _, check := range checks {
		checker, err := checkFactory(check)
		if err != nil {
			fmt.Println(err)
			status = 1
		}
		res, err := checker(logFileDir)
		if err != nil {
			fmt.Println(err)
			status = 1
		}

		res.Output()
	}

	return status
}
