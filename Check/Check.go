package Check
import "errors"

const (
	IMAGE_PULL_BACK_OFF = 0
)

func AllChecks() []int {
	return []int{IMAGE_PULL_BACK_OFF}
}

func GetCheck(checkName int) (Check, error) {
	switch checkName {
	case IMAGE_PULL_BACK_OFF:
		return ImagePullBackOff, nil
	default:
		return nil, errors.New("Could not find requested check")
	}
}
