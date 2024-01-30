package routes

import (
	"fmt"
	"strings"
)

const (
	ActorImages                         string = "/actors/images"
	ActorImageReserved                  string = "/actors/images/reserved"
	ActorImagePathParam                 string = "/actors/images/:id"
	EditActorImagePathParam             string = "/actors/images/:id/edit"
	ActorImageShortDescriptionPathParam string = "/actors/images/:id/sdesc"
	ActorImageDescriptionPathParam      string = "/actors/images/:id/desc"
)

func ActorImagePath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d", ActorImages, id)
	return sb.String()
}

func EditActorImagePath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/edit", ActorImages, id)
	return sb.String()
}

func ActorImageShortDescriptionPath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/sdesc", ActorImages, id)
	return sb.String()
}

func ActorImageDescriptionPath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/desc", ActorImages, id)
	return sb.String()
}
