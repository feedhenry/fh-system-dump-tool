package check
import (
	"io/ioutil"
	"path/filepath"
	"os"
	"encoding/json"
	"strings"
	"fmt"
)


func ImagePullBackOff(logDir string) (CheckResult, error) {
	result := &ImagePullBackOffResult{Status: 0, CheckName: "ImagePullBackOff"}
	projectDirs, err := ioutil.ReadDir(filepath.Join(logDir, "projects"))
	if err != nil {
		return result, err
	}

	for _, projectDir := range projectDirs {
		reader, err := os.Open(filepath.Join(logDir, "projects", projectDir.Name(), "events.json"))
		if err != nil {
			return result, err
		}

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