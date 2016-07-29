package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"fmt"
)

func analyseTask() int {
	res, err := checkForImagePullBackOff(*dumpFileLocation)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	res.Output()

	return 0
}

type checkResult struct {
	Status  int
	Info []checkInfo
}

func (c *checkResult) Output() {
	if c.Status == 0 {
		fmt.Println("Issue not detected.")
		return
	}

	for _, item := range c.Info {
		fmt.Println(item.ObjectName + " has ImagePullBackOff issue")
	}
}

type checkInfo struct {
	File string
	Entry string
	ObjectName string
	Namespace string
	Count int
}


type Events struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
	} `json:"metadata"`
	Items []struct {
		Kind       string `json:"kind"`
		APIVersion string `json:"apiVersion"`
		Metadata   struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			SelfLink          string    `json:"selfLink"`
			UID               string    `json:"uid"`
			ResourceVersion   string    `json:"resourceVersion"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			DeletionTimestamp time.Time `json:"deletionTimestamp"`
		} `json:"metadata"`
		InvolvedObject struct {
			Kind            string `json:"kind"`
			Namespace       string `json:"namespace"`
			Name            string `json:"name"`
			UID             string `json:"uid"`
			APIVersion      string `json:"apiVersion"`
			ResourceVersion string `json:"resourceVersion"`
			FieldPath       string `json:"fieldPath"`
		} `json:"involvedObject"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
		Source  struct {
			Component string `json:"component"`
			Host      string `json:"host"`
		} `json:"source"`
		FirstTimestamp time.Time `json:"firstTimestamp"`
		LastTimestamp  time.Time `json:"lastTimestamp"`
		Count          int       `json:"count"`
		Type           string    `json:"type"`
	} `json:"items"`
}

func checkForImagePullBackOff(logDir string) (*checkResult, error) {
	result := &checkResult{Status: 0}
	projectDirs, err := ioutil.ReadDir(filepath.Join(logDir, "projects"))
	if err != nil {
		return result, err
	}

	for _, projectDir := range projectDirs {
		events := Events{}
		reader, err := os.Open(filepath.Join(logDir, "projects", projectDir.Name(), "events.json"))
		if err != nil {
			return result, err
		}
		decoder := json.NewDecoder(reader)
		decoder.Decode(&events)
		for _, event := range events.Items {
			if event.Reason == "FailedSync" && strings.Contains(event.Message, "ImagePullBackOff") {
				info := checkInfo{ObjectName: event.InvolvedObject.Name, Namespace: event.InvolvedObject.Namespace, Count: event.Count, Entry: event.Message}
				result.Status = 1
				result.Info = append(result.Info, info)
			}
		}
	}

	return result, nil

}
