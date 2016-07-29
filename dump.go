package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// A Task performs some part of the RHMAP System Dump Tool.
type Task func() error

// An errorList accumulates multiple error messages and implements error.
type errorList []string

func (e errorList) Error() string {
	return "multiple errors:\n" + strings.Join(e, "\n")
}


// getSpaceSeparated calls cmd, expected to output a space-separated list of
// words to stdout, and returns the words.
func getSpaceSeparated(cmd *exec.Cmd) ([]string, error) {
	var projects []string
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("command %q: %v", strings.Join(cmd.Args, " "), err)
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		projects = append(projects, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command %q: %v: %v", strings.Join(cmd.Args, " "), err, buf.String())
	}
	return projects, nil
}

func dumpTask() int {
	archiveFile, err := os.Create(*dumpFileLocation)
	if err != nil {
		printError(err)
		return 1
	}
	defer archiveFile.Close()

	tarFile, err := NewTgz(archiveFile)
	if err != nil {
		printError(err)
		return 1
	}
	defer tarFile.Close()

	var tasks []Task

	var resources = []string{"deploymentconfigs", "pods", "services", "events"}

	// Add tasks to fetch resource definitions.
	projects, err := GetProjects()
	if err != nil {
		printError(err)
		return 1
	}

	for _, p := range projects {
		outFor := outToTGZ("json", tarFile)
		errOutFor := outToTGZ("stderr", tarFile)
		task := ResourceDefinitions(p, resources, outFor, errOutFor)
		tasks = append(tasks, task)
	}

	fmt.Println("Starting RHMAP System Dump Tool...")
	defer fmt.Printf("\nDumped system information to: %s\n", archiveFile.Name())

	// Avoid the creating goroutines and other controls if we're executing
	// tasks sequentially.
	if *maxParallelTasks == 1 {
		for _, task := range tasks {
			task()
			fmt.Print(".")
		}
		return 0
	}
	// Run at most N tasks in parallel, and wait for all of them to
	// complete.
	var wg sync.WaitGroup
	sem := make(chan struct{}, *maxParallelTasks)
	for _, task := range tasks {
		task := task
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			task()
			fmt.Print(".")
			<-sem
		}()
	}
	wg.Wait()

	return 0
}

// GetProjects returns a list of project names visible by the current logged in
// user.
func GetProjects() ([]string, error) {
	return getProjects(exec.Command("oc", "get", "projects", "-o=jsonpath={.items[*].metadata.name}"))
}

// getProjects calls cmd, expected to output a space-separated list of project
// names to stdout, and returns a list of project names.
func getProjects(cmd *exec.Cmd) ([]string, error) {
	var projects []string
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("command %q: %v", strings.Join(cmd.Args, " "), err)
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		projects = append(projects, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command %q: %v: %v", strings.Join(cmd.Args, " "), err, buf.String())
	}
	return projects, nil
}
