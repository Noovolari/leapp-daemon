package use_case

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
)

func EditAwsSessionRegion(sessionId string, region string) error {
	/*
		isRegionValid := region2.IsRegionValid(region)
		if !isRegionValid {
			return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + region + " not valid"))
		}

		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		// Find a valid Aws Session
		for _, plainSession := range config.awsIamUserSessions {
			if plainSession.Id == sessionId {
				plainSession.Account.Region = region
				err = configuration.UpdateConfiguration(config, false)
				if err != nil { return err }
				return nil
			}
		}

		for _, federatedSession := range config.AwsIamRoleFederatedSessions {
			if federatedSession.Id == sessionId {
				federatedSession.Account.Region = region
				err = configuration.UpdateConfiguration(config, false)
				if err != nil { return err }
				return nil
			}
		}
	*/

	return http_error.NewUnprocessableEntityError(fmt.Errorf("no valid AWS session found for editing region"))
}
