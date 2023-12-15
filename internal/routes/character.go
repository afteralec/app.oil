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

func CharacterApplicationReviewPath(id string) string {
	return fmt.Sprintf("%s/%s/review", CharacterApplication, id)
}

func SubmitCharacterApplicationPath(id string) string {
	return fmt.Sprintf("%s/%s/submit", CharacterApplication, id)
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

func CharacterApplicationNameReviewPath(id string) string {
	return fmt.Sprintf("%s/%s/name/review", CharacterApplication, id)
}

func CharacterApplicationGenderReviewPath(id string) string {
	return fmt.Sprintf("%s/%s/gender/review", CharacterApplication, id)
}

func CharacterApplicationShortDescriptionReviewPath(id string) string {
	return fmt.Sprintf("%s/%s/sdesc/review", CharacterApplication, id)
}

func CharacterApplicationDescriptionReviewPath(id string) string {
	return fmt.Sprintf("%s/%s/description/review", CharacterApplication, id)
}

func CharacterApplicationBackstoryReviewPath(id string) string {
	return fmt.Sprintf("%s/%s/backstory/review", CharacterApplication, id)
}
