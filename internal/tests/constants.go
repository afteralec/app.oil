package tests

import "fmt"

const (
	TestURL             = "http://petrichormud.com"
	TestUsername        = "testify"
	TestUsernameTwo     = "testify2"
	TestUsernameThree   = "testify3"
	TestPassword        = "T3sted_tested"
	TestEmailAddress    = "testify@test.com"
	TestEmailAddressTwo = "testify2@test.com"
)

func MakeTestURL(path string) string {
	return fmt.Sprintf("%s%s", TestURL, path)
}
