package named_configuration

type NamedConfiguration struct {
	Id   string
	Name string
}

const DefaultNamedConfigurationName = "leapp-default"

func (config *NamedConfiguration) IsDefault() bool {
	return config.Name == DefaultNamedConfigurationName
}
