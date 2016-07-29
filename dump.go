package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// A Task performs some part of the RHMAP System Dump Tool.
type Task func() error

// An errorList accumulates multiple error messages and implements error.
type errorList []string

// A projectResourceWriterFactory generates io.Writers for dumping data of a
// particular resource type within a project.
type projectResourceWriterCloserFactory func(project, resource string) (io.Writer, io.Closer, error)

func (e errorList) Error() string {
	return "multiple errors:\n" + strings.Join(e, "\n")
}

// A getProjectResourceCmdFactory generates commands to get resources of a given
// type in a project.
type getProjectResourceCmdFactory func(project, resource string) *exec.Cmd

// outToFile returns a function that creates an io.Writer that writes to a file
// in basepath with extension, given a project and resource.
func outToFile(basepath, extension string) projectResourceWriterCloserFactory {
	return func(project, resource string) (io.Writer, io.Closer, error) {
		projectpath := filepath.Join(basepath, "projects", project)
		err := os.MkdirAll(projectpath, 0770)
		if err != nil {
			return nil, nil, err
		}
		f, err := os.Create(filepath.Join(projectpath, resource+"."+extension))
		if err != nil {
			return nil, nil, err
		}
		return f, f, nil
	}
}

// OutToTGZ returns an anonymous factory function that will create an io.Writer which writes into the tar archive
// provided. The path inside the tar.gz file is calculated from the project and resource provided
func outToTGZ(extension string, tarFile *Archive) projectResourceWriterCloserFactory {
	return func(project, resource string) (io.Writer, io.Closer, error) {
		projectPath := filepath.Join("projects", project)
		writer := tarFile.GetWriterToFile(filepath.Join(projectPath, resource+"."+extension))
		return writer, writer, nil
	}
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

// ResourceDefinitions fetches the JSON resource definition for all given types
// in project. For each resource type, it uses outFor and errOutFor to get
// io.Writers to write, respectively, the JSON output and any eventual error
// message.
func ResourceDefinitions(project string, types []string, outFor, errOutFor projectResourceWriterCloserFactory) Task {
	return resourceDefinitions(func(project, resource string) *exec.Cmd {
		return exec.Command("oc", "-n", project, "get", resource, "-o=json")
	}, project, types, outFor, errOutFor)
}

func resourceDefinitions(cmdFactory getProjectResourceCmdFactory, project string, types []string, outFor, errOutFor projectResourceWriterCloserFactory) Task {
	return func() error {
		var errors errorList
		// NOTE: we could fetch all resources of all types in a single
		// call to oc, by passing a comma-separated list of resource
		// types. Instead, we call oc multiple times to send the output
		// to different files without processing the contents of the
		// output from oc.
		for _, resource := range types {
			var err error
			var stdoutCloser, stderrCloser io.Closer

			cmd := cmdFactory(project, resource)
			cmd.Stdout, stdoutCloser, err = outFor(project, resource)
			defer stdoutCloser.Close()
			if err != nil {
				errors = append(errors, err.Error())
				// Since we couldn't get an io.Writer for
				// cmd.Stdout, give up processing this resource
				// type, and skip to the next type.
				continue
			}
			var buf bytes.Buffer
			cmd.Stderr, stderrCloser, err = errOutFor(project, resource)
			defer stderrCloser.Close()

			if err != nil {
				errors = append(errors, err.Error())
				// We can possibly try to run the command
				// without an io.Writer from errOutFor. In this
				// case, we'll attach an in-memory buffer so
				// that we can include the stderr output in
				// errors.
				cmd.Stderr = &buf
			} else {
				// Send stderr to both the io.Writer from
				// errOutFor, and an in-memory buffer, used to
				// enrich error messages.
				cmd.Stderr = io.MultiWriter(cmd.Stderr, &buf)
			}
			// TODO: limit the execution time with a timeout.
			err = cmd.Run()
			if err != nil {
				errors = append(errors, fmt.Sprintf("command %q: %v: %v", strings.Join(cmd.Args, " "), err, buf.String()))
				continue
			}
		}
		if len(errors) > 0 {
			return errors
		}
		return nil
	}
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
