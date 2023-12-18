package routes

import "fmt"

const (
	Characters            = "/characters"
	CharacterApplication  = "/character/application"
	CharacterApplications = "/character/applications"
)

func NewCharacterApplicationPath() string {
	return fmt.Sprintf("%s/new", CharacterApplication)
}

func CharacterApplicationPath(id string) string {
	return fmt.Sprintf("%s/%s", CharacterApplication, id)
}

func CharacterApplicationNamePath(id string) string {
	return fmt.Sprintf("%s/%s/name", CharacterApplication, id)
}

func CharacterApplicationGenderPath(id string) string {
	return fmt.Sprintf("%s/%s/gender", CharacterApplication, id)
}

func CharacterApplicationShortDescriptionPath(id string) string {
	return fmt.Sprintf("%s/%s/sdesc", CharacterApplication, id)
}

func CharacterApplicationDescriptionPath(id string) string {
	return fmt.Sprintf("%s/%s/description", CharacterApplication, id)
}

func CharacterApplicationBackstoryPath(id string) string {
	return fmt.Sprintf("%s/%s/backstory", CharacterApplication, id)
}

// TODO: Maybe make these a /status/submitte and /status/in-review, etc route?
func SubmitCharacterApplicationPath(id string) string {
	return fmt.Sprintf("%s/%s/submit", CharacterApplication, id)
}

func PutCharacterApplicationInReviewPath(id string) string {
	return fmt.Sprintf("%s/%s/in-review", CharacterApplication, id)
}

func CharacterApplicationSubmittedPath(id string) string {
	return fmt.Sprintf("%s/%s/submitted", CharacterApplication, id)
}

func CharacterApplicationSubmittedSuccessPath(id string) string {
	return fmt.Sprintf("%s/success", CharacterApplicationSubmittedPath(id))
}

func ReviewCharacterApplicationsPath() string {
	return fmt.Sprintf("%s/review", CharacterApplications)
}
