package check

import (
	"encoding/json"
	"io"
	"strings"
)

type ImagePullBackOff struct {
	Result *Result
}

func NewImagePullBackOff() Check {
	c := &ImagePullBackOff{}
	c.Result = &Result{Status: 0, CheckName: "ImagePullBackOff", StatusMessage: "This issue has not been detected"}
	return c
}

func (c *ImagePullBackOff) GetResult() *Result {
	return c.Result
}

func (c *ImagePullBackOff) RequiredFiles() []string {
	return []string{"events.json"}
}

func (c *ImagePullBackOff) ExamineFile(reader io.Reader) error {

	events := Events{}
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&events)
	if err != nil {
		return err
	}

	for _, event := range events.Items {
		if event.Reason == "FailedSync" && strings.Contains(event.Message, "ImagePullBackOff") {
			info := Info{ObjectName: event.InvolvedObject.Name, Namespace: event.InvolvedObject.Namespace, Count: event.Count, Entry: event.Message}
			c.Result.Status = 1
			c.Result.StatusMessage = "This issue may be present in the system"
			c.Result.Info = append(c.Result.Info, info)
		}
	}

	return nil
}
