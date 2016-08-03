package check

import (
	"errors"
	"io"
)

const (
	IMAGE_PULL_BACK_OFF = 0
)

type CheckFileFactory func([]string) []io.Reader

type CheckTask func(CheckFileFactory) error

func AllChecks() []int {
	return []int{IMAGE_PULL_BACK_OFF}
}

func GetCheck(checkName int) (CheckTask, error) {
	switch checkName {
	case IMAGE_PULL_BACK_OFF:
		return CheckForImagePullBackOff, nil
	default:
		return nil, errors.New("could not find requested check")
	}
}
