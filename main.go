package main

import (
	"leapp_daemon/infrastructure/http/engine"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/providers"
)

func main() {
	//TODO: Move under providers singleton
	defer logging.CloseLogFile()

	prov := providers.NewProviders()
	defer prov.Close()

	config := ConfigurationBootstrap(prov)
	NamedProfilesBootstrap(prov, config)
	AwsIamUserBootstrap(prov, config)
	GcpIamUserAccountOauthBootstrap(prov, config)
	AlibabaRamUserBootstrap(prov, config)
	AlibabaRamRoleFederatedBootstrap(prov, config)
	AlibabaRamRoleChainedBootstrap(prov, config)

	eng := engine.Engine(prov)
	eng.Serve(8080)
}
