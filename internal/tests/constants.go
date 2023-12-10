package tests

import "fmt"

const (
	TestURL         = "http://petrichormud.com"
	TestUsername    = "testify"
	TestUsernameTwo = "testify2"
	TestPassword    = "T3sted_tested"
)

func MakeTestURL(path string) string {
	return fmt.Sprintf("%s%s", TestURL, path)
}
