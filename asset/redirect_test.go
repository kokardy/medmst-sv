package main

import "testing"

func TestDrugCodeToYJ(t *testing.T) {
	drugCode := "225320"
	expected := "4490023F2020"
	yj, err := drugcodeToYJ(drugCode)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if yj != expected {
		t.Fatalf(
			"yj is %s. but expected is %s",
			yj, expected)
	}

	t.Logf(
		"yj is %s. expected is %s",
		yj, expected)

}
