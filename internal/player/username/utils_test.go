package username

import (
	"testing"
)

type SanitizeTestInput struct {
	Input string
	Want  string
}

func TestSanitize(t *testing.T) {
	tests := [3]SanitizeTestInput{{"Test", "test"}, {"t3St", "t3st"}, {"t_st!^&-", "t_st-"}}
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		got := Sanitize(test.Input)
		if got != test.Want {
			t.Errorf("got: %s, wanted: %s", got, test.Want)
		}
	}
}

type isValidTestInput struct {
	Input string
	Want  bool
}

func TestIsValid(t *testing.T) {
	tests := [3]isValidTestInput{{"tst", false}, {"test", true}, {"test-_test-_test-", false}}
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		got := IsValid(test.Input)
		if got != test.Want {
			t.Errorf("got: %t, wanted: %t", got, test.Want)
		}
	}
}
