package Check
import "errors"

const (
	CHECK_IMAGE_PULL_BACK_OFF = 0
	CHECK_PHILS_COOL = 1
)

func GetCheck(checkType int) (Check, error) {
	switch checkType {
	case CHECK_IMAGE_PULL_BACK_OFF:
		return ImagePullBackOff, nil
	default:
		return nil, errors.New("Could not find requested check")
	}
}
