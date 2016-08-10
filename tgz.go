package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"time"
)

type Archive struct {
	tgzFile    io.Writer
	tarWriter  *tar.Writer
	gzWriter   *gzip.Writer
	WriteQueue chan ArchiveWriteTask
	finished   chan struct{}
}

type ArchiveWriter struct {
	File    string
	Archive *Archive
	Writer  *bytes.Buffer
}

type ArchiveWriteTask struct {
	Name    string
	Content []byte
}

func (a *ArchiveWriter) Write(p []byte) (n int, err error) {
	return a.Writer.Write(p)
}

func (a *ArchiveWriter) Close() error {
	// create a write task from the data in this file
	writeTask := ArchiveWriteTask{Name: a.File, Content: a.Writer.Bytes()}
	// and queue it in the archive write queue
	a.Archive.WriteQueue <- writeTask
	return nil
}

func NewTgz(file io.Writer) (*Archive, error) {
	tgz := Archive{}
	var err error
	tgz.tgzFile = file
	if err != nil {
		return nil, err
	}

	tgz.gzWriter = gzip.NewWriter(tgz.tgzFile)
	tgz.tarWriter = tar.NewWriter(tgz.gzWriter)
	tgz.WriteQueue = make(chan ArchiveWriteTask)
	tgz.finished = make(chan struct{})

	go tgz.listenForWrites()

	return &tgz, nil

}

func (a *Archive) listenForWrites() {
	for task := range a.WriteQueue {
		a.AddFileByContent(task.Content, task.Name)
	}
	a.finished <- struct{}{}
}

func (a *Archive) GetWriterToFile(file string) io.WriteCloser {
	writer := ArchiveWriter{File: file, Archive: a, Writer: &bytes.Buffer{}}
	return &writer
}

func (a *Archive) AddFileByContent(src []byte, dest string) error {
	if dest == "" {
		dest = fmt.Sprintf("untitled_%08x", rand.Uint32())
	}

	header := &tar.Header{
		Name:    dest,
		Size:    int64(len(src)),
		Mode:    0775,
		ModTime: time.Now(),
	}

	if err := a.tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(a.tarWriter, bytes.NewReader(src)); err != nil {
		return err
	}

	return nil
}

func (a *Archive) Close() {
	close(a.WriteQueue)

	<-a.finished
	close(a.finished)
	a.tarWriter.Close()
	a.gzWriter.Close()
}
