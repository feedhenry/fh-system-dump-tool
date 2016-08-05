package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"strconv"
	"sync"
	"testing"
)

func TestWritingAFile(t *testing.T) {
	b := bytes.Buffer{}
	tgzFile := bufio.NewReadWriter(bufio.NewReader(&b), bufio.NewWriter(&b))

	tgz, err := NewTgz(tgzFile)
	if err != nil {
		t.Fatal(err)
	}

	writer := tgz.GetWriterToFile("test1.txt")
	writer.Write([]byte("test"))
	writer.Close()
	tgz.Close()
	tgzFile.Flush()

	files, err := decompressAndListFiles(tgzFile)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := files["test1.txt"]; !ok {
		t.Fatal("Expected tgz to contain test1.txt but it didnt")
	}
}

func TestWritingTwoFiles(t *testing.T) {
	b := bytes.Buffer{}
	tgzFile := bufio.NewReadWriter(bufio.NewReader(&b), bufio.NewWriter(&b))

	tgz, err := NewTgz(tgzFile)
	if err != nil {
		t.Fatal(err)
	}

	writer := tgz.GetWriterToFile("test1.txt")
	writer.Write([]byte("test"))
	writer.Close()

	writer = tgz.GetWriterToFile("test2.txt")
	writer.Write([]byte("test"))
	writer.Close()

	tgz.Close()
	tgzFile.Flush()

	files, err := decompressAndListFiles(tgzFile)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := files["test1.txt"]; !ok {
		t.Fatal("Expected tgz to contain test1.txt but it didnt")
	}
	if _, ok := files["test2.txt"]; !ok {
		t.Fatal("Expected tgz to contain test1.txt but it didnt")
	}
}

func TestWritingUTF8ToAFile(t *testing.T) {
	b := bytes.Buffer{}
	tgzFile := bufio.NewReadWriter(bufio.NewReader(&b), bufio.NewWriter(&b))

	tgz, err := NewTgz(tgzFile)
	if err != nil {
		t.Fatal(err)
	}

	writer := tgz.GetWriterToFile("logs/projects/phils-core/pods-fh-aaa-8-v7m10-fh-aaa.logs")
	writer.Write([]byte("世界世形声字 / 形聲字界世界形声字 / 形聲字世界世界世界世形声字 / 形聲字界世界世界世界"))
	writer.Close()
	tgz.Close()
	tgzFile.Flush()

	files, err := decompressAndListFiles(tgzFile)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := files["logs/projects/phils-core/pods-fh-aaa-8-v7m10-fh-aaa.logs"]; !ok {
		t.Fatal("Expected tgz to contain logs/projects/phils-core/pods-fh-aaa-8-v7m10-fh-aaa.logs but it didnt")
	}
}

type task func() error

func getWriteTask(writer io.WriteCloser) task {
	return func() error {
		writer.Write([]byte("世界世形声字 / 形聲字界世界形声字 / 形聲字世界世界世界世形声字 / 形聲字界世界世界世界"))
		writer.Close()
		return nil
	}
}

func TestParallelWrites(t *testing.T) {
	b := bytes.Buffer{}
	tgzFile := bufio.NewReadWriter(bufio.NewReader(&b), bufio.NewWriter(&b))

	tgz, err := NewTgz(tgzFile)
	if err != nil {
		t.Fatal(err)
	}

	numSyncWrites := 10000

	tasks := []task{}

	for i := 1; i <= numSyncWrites; i++ {
		tasks = append(tasks, getWriteTask(tgz.GetWriterToFile("file"+strconv.Itoa(i)+".file")))
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, numSyncWrites)
	for _, task := range tasks {
		task := task
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			task()
			<-sem
		}()
	}
	wg.Wait()

	tgz.Close()
	tgzFile.Flush()

	files, err := decompressAndListFiles(tgzFile)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := files["file1.file"]; !ok {
		t.Fatal("Expected tgz to contain file1.file but it didnt")
	}

	if _, ok := files["file"+strconv.Itoa(numSyncWrites)+".file"]; !ok {
		t.Fatal("Expected tgz to contain file" + strconv.Itoa(numSyncWrites) + ".file but it didnt")
	}

}

func decompressAndListFiles(tgzFile io.Reader) (map[string]string, error) {
	ret := map[string]string{}

	gzf, err := gzip.NewReader(tgzFile)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(gzf)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ret[header.Name] = header.Name
	}

	return ret, nil
}
