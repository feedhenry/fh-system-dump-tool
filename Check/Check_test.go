package Check
import "testing"


func TestAllChecks(t *testing.T) {
	checks := AllChecks()
	if len(checks) < 1 {
		t.Fatal("Count of all checks should be 1 or more in length")
	}

	if checks[0] != IMAGE_PULL_BACK_OFF {
		t.Fatal("first check should be IMAGE_PULL_BACK_OFF")
	}
}