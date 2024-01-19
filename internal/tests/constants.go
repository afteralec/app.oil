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

var TestRoomImage TestRoomImageParams = TestRoomImageParams{
	Name:        "ci-test-room-image",
	Title:       "An elegant, wood-paneled office",
	Description: "Dark, oiled wood encloses this cozy office, each panel polished to an immaculate sheen. In stark contrast, the floor is a pale, sanded expanse of knotted hardwood, with brightly-colored rugs waiting to soften footsteps. A sweeping vista sprawls beyond the floor-to-ceiling windows, its misty landscape dotted with jagged peaks.",
	Size:        "2",
}

func MakeTestURL(path string) string {
	return fmt.Sprintf("%s%s", TestURL, path)
}
