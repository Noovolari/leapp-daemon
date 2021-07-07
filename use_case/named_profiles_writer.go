package use_case

import (
	"leapp_daemon/domain/domain_aws/named_profile"
)

type NamedProfilesWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (namedProfilesWriter *NamedProfilesWriter) UpdateNamedProfiles(oldNamedProfiles []named_profile.NamedProfile, newNamedProfiles []named_profile.NamedProfile) error {
	config, err := namedProfilesWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.NamedProfiles = newNamedProfiles
	err = namedProfilesWriter.ConfigurationRepository.UpdateConfiguration(config)

	return err
}
