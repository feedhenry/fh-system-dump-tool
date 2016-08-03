package check

import (
	"bytes"
	"testing"
)

func TestNewImagePullBackOff(t *testing.T) {
	check := NewImagePullBackOff()
	if check.GetResult().CheckName != "ImagePullBackOff" ||
		check.GetResult().Status != 0 ||
		check.GetResult().StatusMessage != "This issue has not been detected" {
		t.Fatal("check was not constructed correctly")
	}
}

func TestRequiredFiles(t *testing.T) {
	check := NewImagePullBackOff()
	files := check.RequiredFiles()

	expected := []string{"events.json"}
	if len(files) != len(expected) {
		t.Fatalf("expected: %v got: %v", expected, files)
	}
	if files[0] != expected[0] {
		t.Fatalf("expected: %v got: %v", expected, files)
	}
}

func TestDiscoversIssue(t *testing.T) {
	check := NewImagePullBackOff()

	reader := bytes.NewReader([]byte(`
	{
		"items": [
			{
				"involvedObject": {
					"namespace": "phils-core",
					"name": "fh-aaa-4-0w1q2"
				},
				"reason": "FailedSync",
				"message": "Error syncing pod, skipping: failed to \"StartContainer\" for \"fh-aaa\" with ImagePullBackOff: \"Back-off pulling image \\\"docker.io/rhmap/fh-aaa:0.3.0-349-234\\\"\"\n",
				"count": 47
			}
		]
	}
	`))
	err := check.ExamineFile(reader)
	if err != nil {
		t.Fatal(err)
	}

	if check.GetResult().Status != 1 {
		t.Fatalf("check should have set error status, expected 1, got: %v", check.GetResult().Status)
	}
	if check.GetResult().Info[0].Count != 47 {
		t.Fatalf("check should have set count, expected 47, got: %v", check.GetResult().Status)
	}
}

func TestPassesIssues(t *testing.T) {
	check := NewImagePullBackOff()

	reader := bytes.NewReader([]byte(`
	{
		"items": [
			{
				"involvedObject": {
					"namespace": "phils-core",
					"name": "fh-aaa"
				},
				"reason": "DeploymentScaled",
				"message": "Scaled deployment \"fh-aaa-4\" from 1 to 0",
				"count": 1
			}
		]
	}
	`))
	err := check.ExamineFile(reader)
	if err != nil {
		t.Fatal(err)
	}

	if check.GetResult().Status != 0 {
		t.Fatalf("check should have set success status, expected 0, got: %v", check.GetResult().Status)
	}
}

func TestBadJsonIssues(t *testing.T) {
	check := NewImagePullBackOff()

	reader := bytes.NewReader([]byte("{]"))
	err := check.ExamineFile(reader)
	if err == nil {
		t.Fatal("bad JSON should have an error")
	}
}
