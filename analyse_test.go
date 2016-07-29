package main
import (
	"testing"
	"github.com/fheng/fh-system-dump-tool/Check"
	"errors"
)

type Result struct {
	OutputCalled bool
	Dir string
}

func (r *Result) Output() {
	r.OutputCalled = true
}

func TestAllChecksAreExecuted(t *testing.T) {
	checks := []int{0, 1}
	r0 := &Result{OutputCalled: false}
	r1 := &Result{OutputCalled: false}
	checkFactory := func(id int) (Check.Check, error) {
		switch id {
		case 0:
			return func(logDir string) (Check.CheckResult, error) {
				r0.Dir = logDir + "/r0"
				return r0, nil
			}, nil
		case 1:
			return func(logDir string) (Check.CheckResult, error) {
				r1.Dir = logDir + "/r1"
				return r1, nil
			}, nil
		default:
			return nil, errors.New("Could not find id to load")
		}
	}

	ret := analyseTask("foo/bar", checks, checkFactory)

	if ret != 0 {
		t.Fatal("return value was non-zero")
	}

	if r0.OutputCalled != true {
		t.Fatal("result output was not called but should have been")
	}
	if r1.OutputCalled != true {
		t.Fatal("result output was not called but should have been")
	}

	if r0.Dir != "foo/bar/r0" {
		t.Fatal("Check was not passed the correct logDir value got: " + r0.Dir + ", expected: foo/bar/r0")
	}

	if r1.Dir != "foo/bar/r1" {
		t.Fatal("Check was not passed the correct logDir value got: " + r1.Dir + ", expected: foo/bar/r1")
	}

}