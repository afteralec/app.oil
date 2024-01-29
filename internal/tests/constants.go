package tests

import (
	"fmt"

	"petrichormud.com/app/internal/constants"
)

const (
	TestURL             = "http://petrichormud.com"
	TestUsername        = "testify"
	TestUsernameTwo     = "testify2"
	TestUsernameThree   = "testify3"
	TestPassword        = "T3sted_tested"
	TestEmailAddress    = "testify@test.com"
	TestEmailAddressTwo = "testify2@test.com"
)

var TestRoom CreateTestRoomParams = CreateTestRoomParams{
	Title:       "An elegant, wood-paneled office",
	Description: "Dark, oiled wood encloses this cozy office, each panel polished to an immaculate sheen. In stark contrast, the floor is a pale, sanded expanse of knotted hardwood, with brightly-colored rugs waiting to soften footsteps. A sweeping vista sprawls beyond the floor-to-ceiling windows, its misty landscape dotted with jagged peaks.",
	Size:        2,
}

var TestActorImage CreateTestActorImageParams = CreateTestActorImageParams{
	Gender:           constants.GenderObject,
	Name:             "test-actor-image",
	ShortDescription: "glistening handful of pure potential, studded with eyes",
	Description:      "Mucus clings to the subtly-twitching bumps and pocks of this handful of pure potential. Where it runnels into a tear duct or beneath a rubbery eyelid, the eye there blinks - one of many, each with a distinct color and construction. In places it's warm to the touch and others, sickly cold.",
}

func MakeTestURL(path string) string {
	return fmt.Sprintf("%s%s", TestURL, path)
}
