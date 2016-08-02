package check

import (
	"io"
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
	Items []struct {
		InvolvedObject struct {
			Namespace       string `json:"namespace"`
			Name            string `json:"name"`
		} `json:"involvedObject"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
		Count          int       `json:"count"`
	} `json:"items"`
}
