package check

import (
	"io"
	"time"
)

type Factory func(checkName int) (Check, error)

type Info struct {
	File       string
	Entry      string
	ObjectName string
	Namespace  string
	Count      int
}

type Result struct {
	CheckName     string `json:"checkName" yaml:"checkName"`
	Status        int    `json:"status" yaml:"status"`
	StatusMessage string `json:"statusMessage" yaml:"statusMessage"`
	Info          []Info `json:"info" yaml:"info"`
}

type Check interface {
	ExamineFile(file io.Reader) error
	RequiredFiles() []string
	GetResult() *Result
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
