package use_case

import (
	"leapp_daemon/domain/domain_aws/named_profile"
)

type NamedProfilesActions struct {
	Environment         Environment
	NamedProfilesFacade NamedProfilesFacade
}

func (actions *NamedProfilesActions) GetNamedProfiles() []named_profile.NamedProfile {
	return actions.NamedProfilesFacade.GetNamedProfiles()
}

func (actions *NamedProfilesActions) GetNamedProfileById(profileId string) (named_profile.NamedProfile, error) {
	return actions.NamedProfilesFacade.GetNamedProfileById(profileId)
}

func (actions *NamedProfilesActions) GetOrCreateNamedProfile(profileName string) (named_profile.NamedProfile, error) {
	if profileName == "" {
		profileName = named_profile.DefaultNamedProfileName
	}

	facade := actions.NamedProfilesFacade
	namedProfile, err := facade.GetNamedProfileByName(profileName)
	if err != nil {
		namedProfile = named_profile.NamedProfile{
			Id:   actions.Environment.GenerateUuid(),
			Name: profileName,
		}
		err = facade.AddNamedProfile(namedProfile)
		if err != nil {
			return named_profile.NamedProfile{}, err
		}
	}

	return namedProfile, nil
}
