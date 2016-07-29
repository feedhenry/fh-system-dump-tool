// The RHMAP System Dump Tool.
package main

import (
	"flag"
	"fmt"
	"github.com/fheng/fh-system-dump-tool/check"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	// dumpDir is a path to the base directory where the output of the tool
	// goes.
	dumpDir = "rhmap-dumps"
	// dumpTimestampFormat is a layout for use with Time.Format. Used to
	// create directories with a timestamp. Based on time.RFC3339.
	dumpTimestampFormat = "2006-01-02T15-04-05Z0700"
)

var checks = check.AllChecks()

var maxParallelTasks = flag.Int("p", runtime.NumCPU(), "max number of tasks to run in parallel")
var dumpFileLocation = flag.String("f", "", "The location of the dump file")
var taskType = flag.String("t", "help", "The task to execute")

func printError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
}

func main() {
	flag.Parse()
	if !(*maxParallelTasks > 0) {
		printError(fmt.Errorf("argument to -p flag must be greater than 0"))
		os.Exit(1)
	}

	if *dumpFileLocation == "" {
		start := time.Now().UTC()
		startTimestamp := start.Format(dumpTimestampFormat)
		*dumpFileLocation = filepath.Join(dumpDir, startTimestamp) + ".tar.gz"
	}

	switch *taskType {
	case "dump":
		os.Exit(dumpTask())
	case "analyse":
		os.Exit(analyseTask(*dumpFileLocation, checks, check.GetCheck))
	case "help":
		fallthrough
	default:
		os.Exit(helpTask())
	}
}
