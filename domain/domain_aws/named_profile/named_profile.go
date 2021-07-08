package named_profile

type NamedProfile struct {
	Id   string
	Name string
}

const DefaultNamedProfileName = "default"

func (prof *NamedProfile) IsDefault() bool {
	return prof.Name == DefaultNamedProfileName
}
