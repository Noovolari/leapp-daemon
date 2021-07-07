package domain

import (
	"encoding/json"
	"leapp_daemon/domain/domain_aws/aws_iam_role_chained"
	"leapp_daemon/domain/domain_aws/aws_iam_role_federated"
	"leapp_daemon/domain/domain_aws/aws_iam_user"
	"leapp_daemon/domain/domain_aws/named_profile"
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
)

type Configuration struct {
	ProxyConfiguration             ProxyConfiguration
	AwsIamUserSessions             []aws_iam_user.AwsIamUserSession
	AwsIamRoleChainedSessions      []aws_iam_role_chained.AwsIamRoleChainedSession
	AwsIamRoleFederatedSessions    []aws_iam_role_federated.AwsIamRoleFederatedSession
	GcpIamUserAccountOauthSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession
	NamedProfiles                  []named_profile.NamedProfile
}

type ProxyConfiguration struct {
	ProxyProtocol string
	ProxyUrl      string
	ProxyPort     uint64
	Username      string
	Password      string
}

func GetDefaultConfiguration() Configuration {
	return Configuration{
		ProxyConfiguration: ProxyConfiguration{
			ProxyProtocol: "https",
			ProxyUrl:      "",
			ProxyPort:     8080,
			Username:      "",
			Password:      "",
		},
		AwsIamUserSessions:             make([]aws_iam_user.AwsIamUserSession, 0),
		AwsIamRoleChainedSessions:      make([]aws_iam_role_chained.AwsIamRoleChainedSession, 0),
		AwsIamRoleFederatedSessions:    make([]aws_iam_role_federated.AwsIamRoleFederatedSession, 0),
		GcpIamUserAccountOauthSessions: make([]gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession, 0),
		NamedProfiles:                  make([]named_profile.NamedProfile, 0),
	}
}

func FromJson(configurationJson string) Configuration {
	var config Configuration
	_ = json.Unmarshal([]byte(configurationJson), &config)
	return config
}
