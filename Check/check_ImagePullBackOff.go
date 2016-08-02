package check
import (
	"encoding/json"
	"strings"
	"fmt"
	"io"
)

type ImagePullBackOff struct {

}

func (c *ImagePullBackOff) RequiredFiles() []string {
	return []string{"events.json"}
}

func (c *ImagePullBackOff) Execute(files []io.Reader) (CheckResult, error) {
	result := &ImagePullBackOffResult{Status: 0, CheckName: "ImagePullBackOff"}

	for _, reader := range files {
		events := Events{}
		decoder := json.NewDecoder(reader)
		decoder.Decode(&events)

		for _, event := range events.Items {
			if event.Reason == "FailedSync" && strings.Contains(event.Message, "ImagePullBackOff") {
				info := Info{ObjectName: event.InvolvedObject.Name, Namespace: event.InvolvedObject.Namespace, Count: event.Count, Entry: event.Message}
				result.Status = 1
				result.Info = append(result.Info, info)
			}
		}
	}

	return result, nil
}


type ImagePullBackOffResult struct {
	Status    int
	CheckName string
	Info      []Info
}

func (c *ImagePullBackOffResult) Output() {
	fmt.Println(c.CheckName + " results: ")
	if c.Status == 0 {
		fmt.Println("	✔ - Issue not detected.")
		return
	}

	projectData := map[string]map[string]string{}

	for _, item := range c.Info {
		if _, ok := projectData[item.Namespace]; !ok {
			projectData[item.Namespace] = map[string]string{}
		}
		projectData[item.Namespace][item.ObjectName] = item.Entry
	}

	for projectName, project := range projectData {
		fmt.Println("	Project: " + projectName)
		for podName, msg := range project {
			fmt.Println("		Pod: " + podName)
			fmt.Println("			Msg: " + msg)
		}
	}
}