package check

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

func TestCheckFactory(t *testing.T) {
	_, err := GetCheck(-1)
	if err == nil {
		t.Fatal("Factory should return an error when a bad id is requested")
	}

	checker, err := GetCheck(IMAGE_PULL_BACK_OFF)
	if err != nil {
		t.Fatal(err)
	}

	if checker == nil {
		t.Fatal("Check should not be nil")
	}
}