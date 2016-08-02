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
		t.Fatal("Check was not constructed correctly")
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
		"kind": "List",
		"apiVersion": "v1",
		"metadata": {},
		"items": [
			{
				"kind": "Event",
				"apiVersion": "v1",
				"metadata": {
					"name": "fh-aaa-4-0w1q2.14656378b8b4a890",
					"namespace": "phils-core",
					"selfLink": "/api/v1/namespaces/phils-core/events/fh-aaa-4-0w1q2.14656378b8b4a890",
					"uid": "488dcbb8-5493-11e6-b504-0800275732c8",
					"resourceVersion": "47238",
					"creationTimestamp": "2016-07-28T07:17:03Z",
					"deletionTimestamp": "2016-07-28T09:28:37Z"
				},
				"involvedObject": {
					"kind": "Pod",
					"namespace": "phils-core",
					"name": "fh-aaa-4-0w1q2",
					"uid": "4539b140-5493-11e6-b504-0800275732c8",
					"apiVersion": "v1",
					"resourceVersion": "46966"
				},
				"reason": "FailedSync",
				"message": "Error syncing pod, skipping: failed to \"StartContainer\" for \"fh-aaa\" with ImagePullBackOff: \"Back-off pulling image \\\"docker.io/rhmap/fh-aaa:0.3.0-349-234\\\"\"\n",
				"source": {
					"component": "kubelet",
					"host": "local.feedhenry.io"
				},
				"firstTimestamp": "2016-07-28T07:17:03Z",
				"lastTimestamp": "2016-07-28T07:28:37Z",
				"count": 47,
				"type": "Warning"
			}
		]
	}
	`))
	err := check.ExamineFile(reader)
	if err != nil {
		t.Fatal(err)
	}

	if check.GetResult().Status != 1 {
		t.Fatalf("Check should have set error status, expected 1, got: %v", check.GetResult().Status)
	}
	if check.GetResult().Info[0].Count != 47 {
		t.Fatalf("Check should have set count, expected 47, got: %v", check.GetResult().Status)
	}
}

func TestPassesIssues(t *testing.T) {
	check := NewImagePullBackOff()

	reader := bytes.NewReader([]byte(`
	{
		"kind": "List",
		"apiVersion": "v1",
		"metadata": {},
		"items": [
			{
				"kind": "Event",
				"apiVersion": "v1",
				"metadata": {
					"name": "fh-aaa.1465641d2031006e",
					"namespace": "phils-core",
					"selfLink": "/api/v1/namespaces/phils-core/events/fh-aaa.1465641d2031006e",
					"uid": "ed6dc031-5494-11e6-b504-0800275732c8",
					"resourceVersion": "47242",
					"creationTimestamp": "2016-07-28T07:28:49Z",
					"deletionTimestamp": "2016-07-28T09:28:49Z"
				},
				"involvedObject": {
					"kind": "DeploymentConfig",
					"namespace": "phils-core",
					"name": "fh-aaa",
					"uid": "772cbd90-524f-11e6-b504-0800275732c8",
					"apiVersion": "v1",
					"resourceVersion": "46946"
				},
				"reason": "DeploymentScaled",
				"message": "Scaled deployment \"fh-aaa-4\" from 1 to 0",
				"source": {
					"component": "deploymentconfig-controller"
				},
				"firstTimestamp": "2016-07-28T07:28:49Z",
				"lastTimestamp": "2016-07-28T07:28:49Z",
				"count": 1,
				"type": "Normal"
			}
		]
	}
	`))
	err := check.ExamineFile(reader)
	if err != nil {
		t.Fatal(err)
	}

	if check.GetResult().Status != 0 {
		t.Fatalf("Check should have set success status, expected 0, got: %v", check.GetResult().Status)
	}
}

func TestBadJsonIssues(t *testing.T) {
	check := NewImagePullBackOff()

	reader := bytes.NewReader([]byte("{]"))
	err := check.ExamineFile(reader)
	if err == nil {
		t.Fatal("Bad JSON should have an error")
	}
}
