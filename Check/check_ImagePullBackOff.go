package check

import (
	"encoding/json"
	"io"
	"strings"
)

func CheckForImagePullBackOff(fileFactory CheckFileFactory) Result {
	files := fileFactory([]string{"events.json"})

	res := Result{Status: 0, CheckName: "ImagePullBackOff", StatusMessage: "This issue has not been detected"}

	for _, reader := range files {
		events := Events{}
		decoder := json.NewDecoder(reader)
		err := decoder.Decode(&events)
		if err != nil {
			return err
		}


		for _, event := range events.Items {
			if event.Reason == "FailedSync" && strings.Contains(event.Message, "ImagePullBackOff") {
				info := Info{ObjectName: event.InvolvedObject.Name, Namespace: event.InvolvedObject.Namespace, Count: event.Count, Entry: event.Message}
				res.Status = 1
				res.StatusMessage = "This issue may be present in the system"
				res.Info = append(res.Info, info)
			}
		}
	}

	return res
}
