package request

const (
	FieldName             = "name"
	FieldGender           = "gender"
	FieldShortDescription = "sdesc"
	FieldDescription      = "description"
	FieldBackstory        = "backstory"
)

var FieldsByType map[string]map[string]bool = map[string]map[string]bool{
	TypeCharacterApplication: {
		FieldName:             true,
		FieldGender:           true,
		FieldShortDescription: true,
		FieldDescription:      true,
		FieldBackstory:        true,
	},
}
