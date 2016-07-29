package main

import (
	"github.com/fheng/fh-system-dump-tool/Check"
	"fmt"
)

func analyseTask() int {
	checks := []int{Check.CHECK_IMAGE_PULL_BACK_OFF}



	for _, check := range checks {
		checker, err := Check.GetCheck(check)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		res, err := checker(*dumpFileLocation)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		res.Output()

	}

	return 0
}

