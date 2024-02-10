package routes

import (
	"fmt"
	"strings"
)

const (
	ActorImages                         string = "/actor/images"
	ActorImageReserved                  string = "/actor/images/reserved"
	ActorImagePathParam                 string = "/actor/images/:id"
	EditActorImagePathParam             string = "/actor/images/:id/edit"
	ActorImageShortDescriptionPathParam string = "/actor/images/:id/sdesc"
	ActorImageDescriptionPathParam      string = "/actor/images/:id/desc"
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
