// The RHMAP System Dump Tool.
package main

import (
	"flag"
	"fmt"
	"github.com/fheng/fh-system-dump-tool/check"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
var outputFormat = flag.String("o", "json", "The format to output analysis results in, only JSON currently supported")
var prettyOutputFormat = flag.Bool("pretty", false, "Whether to pretty print the JSON output or not")
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

	outputFormatString := strings.ToLower(*outputFormat)
	if outputFormatString != "yaml" && outputFormatString != "json" {
		printError(fmt.Errorf("argument to -o flag must be either yaml or json"))
		os.Exit(1)
	}

	switch *taskType {
	case "dump":
		err := dumpTask()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "analyse":

		//todo: extract tar file into a dump directory

		files, err := listAllFiles(*dumpFileLocation)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		results, err := analyseTask(files, checks, check.GetCheck)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		outputResults(results, outputFormatString, *prettyOutputFormat)

	case "help":
		fallthrough
	default:
		help, err := ioutil.ReadFile("./help.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print(string(help))
	}
}
