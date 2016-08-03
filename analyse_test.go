package main

import (
	"errors"
	"github.com/fheng/fh-system-dump-tool/check"
	"io"
	"strconv"
	"testing"
)

type mockCheck struct {
	Result *check.Result
	Files  []string
}

func (m *mockCheck) ExamineFile(reader io.Reader) error {
	return nil
}

func (m *mockCheck) RequiredFiles() []string {
	return m.Files
}

func (m *mockCheck) GetResult() *check.Result {
	return m.Result
}

func checkFactory(id int) (check.Check, error) {
	result := &check.Result{StatusMessage: "no issues", Status: id, CheckName: "test check " + strconv.Itoa(id)}
	switch id {
	case 0, 1:
		return &mockCheck{Result: result}, nil
	default:
		return nil, errors.New("invalid check specified")
	}
}

func TestAllChecksReturnCorrectly(t *testing.T) {
	checks := []int{0, 1}

	results, err := analyseTask([]string{}, checks, checkFactory)

	if err != nil {
		t.Fatal(err)
	}

	r0 := results["results"][0]
	r1 := results["results"][1]
	if r0.CheckName != "test check 0" {
		t.Fatal("first test name not returned correctly")
	}
	if r0.Status != 0 {
		t.Fatal("first status not return correctly")
	}

	if r1.CheckName != "test check 1" {
		t.Fatal("second test name not returned correctly")
	}
	if r1.Status != 1 {
		t.Fatal("second status not return correctly")
	}
}

func TestMatchingFiles(t *testing.T) {
	check := &mockCheck{Files: []string{"test1.json", "test.txt"}}
	files := []string{"/path/to/test1.json", "test.txt", "/bad/file/test1.json.txt", "/path/to/test.txt"}
	ret, err := getFilesFor(check, files)
	if err != nil {
		t.Fatal(err)
	}

	if ret[0] != "/path/to/test1.json" {
		t.Fatal("missing expected file: /path/to/test1.json")
	}

	if ret[1] != "test.txt" {
		t.Fatal("missing expected file: test.txt")
	}

	if ret[2] != "/path/to/test.txt" {
		t.Fatal("missing expected file: /path/to/test.txt")
	}

}
