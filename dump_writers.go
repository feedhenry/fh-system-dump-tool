package main

import (
	"io"
	"os"
	"path/filepath"
)

// A projectResourceWriterFactory generates io.Writers for dumping data of a
// particular resource type within a project.
type projectResourceWriterCloserFactory func(project, resource string) (io.Writer, io.Closer, error)

// A platformResourceWriterFactory generates io.Writers for dumping data of a
// particular resource type.
type platformResourceWriterCloserFactory func(resource string) (io.Writer, io.Closer, error)

// outToFileForProject returns a function that creates an io.Writer that writes to a file
// in basepath with extension, given a project and resource. The scope can be passed as a value to group the
// files under, or assigned an empty value to have the files grouped directly under the project.
func outToFileForProject(basepath, extension, scope string) projectResourceWriterCloserFactory {
	return func(project, resource string) (io.Writer, io.Closer, error) {
		scopePath := filepath.Join(basepath, "projects", project, scope)
		err := os.MkdirAll(scopePath, 0770)
		if err != nil {
			return nil, nil, err
		}
		f, err := os.Create(filepath.Join(scopePath, resource+"."+extension))
		if err != nil {
			return nil, nil, err
		}
		return f, f, nil
	}
}

// outToFileForPlatform returns a function that creates an io.Writer that writes to a file
// in basepath with extension, given a project and resource. The scope can be passed as a value to group the
// files under, or assigned an empty value to have the files grouped directly under the project.
func outToFileForPlatform(basepath, extension, scope string) platformResourceWriterCloserFactory {
	return func(resource string) (io.Writer, io.Closer, error) {
		scopePath := filepath.Join(basepath, scope)
		err := os.MkdirAll(scopePath, 0770)
		if err != nil {
			return nil, nil, err
		}
		f, err := os.Create(filepath.Join(scopePath, resource+"."+extension))
		if err != nil {
			return nil, nil, err
		}
		return f, f, nil
	}
}
