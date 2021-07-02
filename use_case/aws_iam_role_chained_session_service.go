package use_case

import (
	"fmt"
	"leapp_daemon/domain/configuration"
	session2 "leapp_daemon/domain/session"
	http_error2 "leapp_daemon/infrastructure/http/http_error"
)

func CreateAwsIamRoleChainedSession(parentId string, accountName string, accountNumber string, roleName string, region string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = CheckParentExist(parentId, config)
		if err != nil {
			return err
		}

		sessions, err := config.GetAwsIamRoleChainedSessions()
		if err != nil {
			return err
		}

		for _, sess := range sessions {
			account := sess.Account
			if sess.ParentId == parentId && account.AccountNumber == accountNumber && account.Role.Name == roleName {
				err := http_error2.NewConflictError(fmt.Errorf("a session with the same parent, account number and role name already exists"))
				return err
			}
		}

		trustedAwsAccount := session2.AwsIamRoleChainedAccount{
			AccountNumber: accountNumber,
			Name:          accountName,
			Role: &session2.AwsIamRole{
				Name: roleName,
				Arn:  fmt.Sprintf("arn:aws:iam::%s:role/%s", accountNumber, roleName),
			},
			Region: region,
		}

		// TODO check uuid format
		uuidString := uuid.New().String() //use Environment.GenerateUuid()
		uuidString = strings.Replace(uuidString, "-", "", -1)

		sess := session2.AwsIamRoleChainedSession{
			Id:        uuidString,
			Status:    session2.NotActive,
			StartTime: "",
			ParentId:  parentId,
			Account:   &trustedAwsAccount,
		}

		err = config.SetAwsIamRoleChainedSessions(append(sessions, &sess))
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func GetAwsIamRoleChainedSession(id string) (*session2.AwsIamRoleChainedSession, error) {
	/*
		var sess *session2.AwsIamRoleChainedSession

		config, err := configuration.ReadConfiguration()
		if err != nil {
			return sess, err
		}

		sessions, err := config.GetAwsIamRoleChainedSessions()
		if err != nil {
			return nil, err
		}

		for _, s := range sessions {
			if s.Id == id {
				return s, nil
			}
		}
	*/

	return nil, http_error2.NewNotFoundError(fmt.Errorf("no session found with id %s", id))
}

func UpdateAwsIamRoleChainedSession(id string, parentId string, accountName string, accountNumber string, roleName string, region string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		foundId := false
		sessions, err := config.GetAwsIamRoleChainedSessions()
		if err != nil {
			return err
		}

		for _, s := range sessions {
			if s.Id == id {
				foundId = true

				if parentId != "" {
					err := CheckParentExist(parentId, config)
					if err != nil {
						return err
					}
					s.ParentId = parentId
				}

				if accountNumber != "" {
					s.Account.AccountNumber = accountNumber
				}

				if accountName != "" {
					s.Account.Name = accountName
				}
				if roleName != "" {
					s.Account.Role.Name = roleName
				}
				if region != "" {
					s.Account.Region = region
				}

				if accountNumber != "" || roleName != "" {
					s.Account.Role.Arn = fmt.Sprintf("arn:aws:iam::%s:role/%s", s.Account.AccountNumber, s.Account.Role.Name)
				}

				break
			}
		}

		if !foundId {
			return http_error2.NewNotFoundError(fmt.Errorf("no session found with id %s", id))
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func DeleteAwsIamRoleChainedSession(id string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		sessions, err := config.GetAwsIamRoleChainedSessions()
		if err != nil {
			return err
		}

		found := false
		for index := range sessions {
			if sessions[index].Id == id {
				sessions = append(sessions[:index], sessions[index+1:]...)
				found = true
				break
			}
		}

		if found == false {
			err = http_error2.NewNotFoundError(fmt.Errorf("trusted aws session with id %s not found", id))
			return err
		}

		err = config.SetAwsIamRoleChainedSessions(sessions)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func CheckParentExist(parentId string, config *configuration.Configuration) error {
	/*
		foundId := false
		plains, err := config.GetSessions()
		if err != nil {
			return err
		}
		for _, sess := range plains {
			if sess.Id == parentId {
				foundId = true
			}
		}

		feds, err := config.AwsGetIamRoleFederatedSessions()
		if err != nil {
			return err
		}
		for _, sess := range feds {
			if sess.Id == parentId {
				foundId = true
			}
		}

		if !foundId {
			err := http_error2.NewNotFoundError(fmt.Errorf("no plain or federated session with id %s found", parentId))
			return err
		}
	*/

	return nil
}
