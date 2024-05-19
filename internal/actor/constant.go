package actor

import (
	"petrichormud.com/app/internal/util"
)

const (
	GenderMale      = "Male"
	GenderFemale    = "Female"
	GenderNonBinary = "NonBinary"
	GenderObject    = "Object"
)

const (
	DefaultCharacterName         string = "Unnamed"
	DefaultImageDescription      string = "Mucus clings to the subtly-twitching bumps and pocks of this handful of pure potential. Where it runnels into a tear duct or beneath a rubbery eyelid, the eye there blinks - one of many, each with a distinct color and construction. In places it's warm to the touch and others, sickly cold."
	DefaultImageShortDescription string = "glistening handful of pure potential, studded with eyes"
	DefaultImageGender           string = GenderObject
)

const (
	ImageNameMinLen        int    = 4
	ImageNameMaxLen        int    = 50
	ImageNameRegex         string = "[^a-z-]+"
	ShortDescriptionMinLen int    = 8
	ShortDescriptionMaxLen int    = 300
	ShortDescriptionRegex  string = "[^a-zA-Z, -]+"
	DescriptionMinLen      int    = 32
	DescriptionMaxLen      int    = 2000
	DescriptionRegex       string = "[^a-zA-Z, '-.!()]+"
)

const (
	CharacterNameMinLen      int    = 4
	CharacterNameMaxLen      int    = 16
	CharacterNameRegex       string = "[^a-zA-Z'-]+"
	CharacterBackstoryMinLen int    = 500
	CharacterBackstoryMaxLen int    = 10000
	CharacterBackstoryRegex  string = "[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+"
)

var (
	GenderMinLen int    = util.MinLengthOfStrings([]string{GenderNonBinary, GenderFemale, GenderMale})
	GenderMaxLen int    = util.MaxLengthOfStrings([]string{GenderNonBinary, GenderFemale, GenderMale})
	GenderRegex  string = util.RegexForExactMatchStrings([]string{GenderNonBinary, GenderFemale, GenderMale})
)

var (
	KeywordMinLen int    = 2
	KeywordMaxLen int    = ShortDescriptionMaxLen
	KeywordRegex  string = "[^a-zA-Z]+"
)
