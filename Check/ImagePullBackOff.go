package Check
import (
	"io/ioutil"
	"path/filepath"
	"os"
	"encoding/json"
	"strings"
)


func ImagePullBackOff(logDir string) (*Result, error) {
	result := &Result{Status: 0, CheckName: "ImagePullBackOff"}
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

