package main

import (
	"fmt"
	"github.com/fheng/fh-system-dump-tool/Check"
)

func analyseTask(checks []int) int {
	status := 0
	for _, check := range checks {
		checker, err := Check.GetCheck(check)
		if err != nil {
			fmt.Println(err)
			status = 1
		}
		res, err := checker(*dumpFileLocation)
		if err != nil {
			fmt.Println(err)
			status = 1
		}

		res.Output()
	}

	return status
}
