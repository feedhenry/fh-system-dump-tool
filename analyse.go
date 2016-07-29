package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"strconv"
)

func analyseTask() int {
	checkForImagePullBackOff(*dumpFileLocation)
	return 0
}

type checkResult struct {
	Status  int
	Message string
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

func checkForImagePullBackOff(logDir string) error {

	projectDirs, err := ioutil.ReadDir(filepath.Join(logDir, "projects"))
	if err != nil {
		return err
	}

	for _, projectDir := range projectDirs {
		events := Events{}
		reader, err := os.Open(filepath.Join(logDir, "projects", projectDir.Name(), "events.json"))
		if err != nil {
			return err
		}
		decoder := json.NewDecoder(reader)
		decoder.Decode(&events)
		for _, event := range events.Items {
			if event.Reason == "FailedSync" && strings.Contains(event.Message, "ImagePullBackOff") {
				fmt.Println(event.InvolvedObject.Name + " in project " + event.InvolvedObject.Namespace + " has failed to pull the image " + strconv.Itoa(event.Count) + "times.")
				fmt.Println("Original Message: " + event.Message)
			}
		}
	}

	return nil

}
