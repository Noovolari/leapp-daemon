package use_case

func RotateAllSessionsCredentials() error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil { return err }

		sessions := config.GetAllSessions()

		for _, sess := range sessions {
	    err = sess.Rotate(&session2.RotateConfiguration{})
	    if err != nil {
	      return err
	    }
	  }

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func RotateSessionCredentialsWithMfaToken(sessionId string, mfaToken string) error {
	/*
		found := false

		config, err := configuration.ReadConfiguration()
		if err != nil { return err }

		for i := range config.awsIamUserSessions {
			sess := config.awsIamUserSessions[i]
			if sess.ID == sessionId {
				err = sess.Rotate(&session2.RotateConfiguration{MfaToken: mfaToken})
				if err != nil {
					return err
				}
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("No session found with id " + sessionId)
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}
