package tests

import (
	"fmt"

	"petrichormud.com/app/internal/actor"
	"petrichormud.com/app/internal/constant"
)

const (
	TestURL             = "http://petrichormud.com"
	TestUsername        = "testify"
	TestUsernameTwo     = "testify2"
	TestUsernameThree   = "testify3"
	TestPassword        = "T3sted_tested"
	TestEmailAddress    = "testify@test.com"
	TestEmailAddressTwo = "testify2@test.com"
	TestActorImageName  = "test-actor-image"
)

var TestRoom CreateTestRoomParams = CreateTestRoomParams{
	Title:       "An elegant, wood-paneled office",
	Description: "Dark, oiled wood encloses this cozy office, each panel polished to an immaculate sheen. In stark contrast, the floor is a pale, sanded expanse of knotted hardwood, with brightly-colored rugs waiting to soften footsteps. A sweeping vista sprawls beyond the floor-to-ceiling windows, its misty landscape dotted with jagged peaks.",
	Size:        2,
}

var TestActorImage CreateTestActorImageParams = CreateTestActorImageParams{
	Gender:           constant.GenderObject,
	Name:             TestActorImageName,
	ShortDescription: actor.DefaultImageShortDescription,
	Description:      actor.DefaultImageDescription,
}

var TestHelpFile CreateTestHelpFileParams = CreateTestHelpFileParams{
	Slug:     "test",
	Title:    "Test Help File",
	Sub:      "A test help file",
	Category: "Test",
	Raw:      "# Test",
	HTML:     "<h1>Test</h1>",
	Related:  []CreateTestHelpFileRelatedParams{},
	Tags:     []CreateTestHelpFileTagParams{},
	PID:      0,
}

func MakeTestURL(path string) string {
	return fmt.Sprintf("%s%s", TestURL, path)
}
