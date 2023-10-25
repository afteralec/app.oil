package password

import (
	"strings"
	"testing"
)

type ValidateTestInput struct {
	Input string
	Want  bool
}

func TestValidate(t *testing.T) {
	tests := [4]ValidateTestInput{{strings.Repeat("a", 3), false}, {strings.Repeat("a", 8), true}, {strings.Repeat("a", 255), true}, {strings.Repeat("a", 256), false}}
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		got := Validate(test.Input)
		if got != test.Want {
			t.Errorf("got: %t, wanted: %t", got, test.Want)
		}
	}
}
