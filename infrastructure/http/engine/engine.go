package engine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"leapp_daemon/infrastructure/http/middleware"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http_controller"
	"leapp_daemon/providers"
)

type engineWrapper struct {
	providers *providers.Providers
	ginEngine *gin.Engine
}

var engineWrapperInstance *engineWrapper = nil

func newEngineWrapper(providers *providers.Providers) *engineWrapper {
	ginEngine := gin.New()

	engineWrapper := engineWrapper{
		ginEngine: ginEngine,
		providers: providers,
	}

	engineWrapper.initialize()

	return &engineWrapper
}

func Engine(providers *providers.Providers) *engineWrapper {
	if engineWrapperInstance != nil {
		return engineWrapperInstance
	} else {
		engineWrapperInstance = newEngineWrapper(providers)
		return engineWrapperInstance
	}
}

func (engineWrapper *engineWrapper) initialize() {
	logging.InitializeLogger()
	engineWrapper.ginEngine.Use(middleware.ErrorHandler.Handle)
	initializeRoutes(engineWrapper.ginEngine, engineWrapper.providers)
}

func (engineWrapper *engineWrapper) Serve(port int) {
	err := engineWrapper.ginEngine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalln("error:", err.Error())
	}
}

func initializeRoutes(ginEngine *gin.Engine, providers *providers.Providers) {
	contr := http_controller.EngineController{Providers: providers}

	v1 := ginEngine.Group("/api/v1")
	{
		// General
		v1.GET("sessions", contr.ListSession)

		// AWS sessions
		v1.GET("aws/named-profiles", contr.GetNamedProfiles)
		v1.GET("aws/regions", contr.GetAwsRegionList)
		v1.PUT("aws/sessions/:id/region", contr.EditAwsRegion)

		// AWS IAM UserName sessions
		v1.POST("aws/iam-user-sessions", contr.CreateAwsIamUserSession)
		v1.PUT("aws/iam-user-sessions/:id", contr.EditAwsIamUserSession)
		v1.GET("aws/iam-user-sessions/:id", contr.GetAwsIamUserSession)
		v1.DELETE("aws/iam-user-sessions/:id", contr.DeleteAwsIamUserSession)
		v1.POST("aws/iam-user-sessions/:id/start", contr.StartAwsIamUserSession)
		v1.POST("aws/iam-user-sessions/:id/stop", contr.StopAwsIamUserSession)
		v1.POST("aws/iam-user-sessions/:id/confirm-mfa-token", contr.ConfirmMfaToken)

		// AWS IAM Role Federated sessions
		v1.GET("aws/iam-role-federated-sessions/:id", contr.GetAwsIamRoleFederatedSession)
		v1.POST("aws/iam-role-federated-sessions", contr.CreateAwsIamRoleFederatedSession)
		v1.PUT("aws/iam-role-federated-sessions/:id", contr.EditAwsIamRoleFederatedSession)
		v1.DELETE("aws/iam-role-federated-sessions/:id", contr.DeleteAwsIamRoleFederatedSession)
		v1.POST("aws/iam-role-federated-sessions/:id/start", contr.StartAwsIamRoleFederatedSession)
		v1.POST("aws/iam-role-federated-sessions/:id/stop", contr.StopAwsIamRoleFederatedSession)

		// AWS IAM Role Chained sessions
		v1.GET("aws/iam-role-chained-sessions/:id", contr.GetAwsIamRoleChainedSession)
		v1.POST("aws/iam-role-chained-sessions", contr.CreateAwsIamRoleChainedSession)
		v1.PUT("aws/iam-role-chained-sessions/:id", contr.EditAwsIamRoleChainedSession)
		v1.DELETE("aws/iam-role-chained-sessions/:id", contr.DeleteAwsIamRoleChainedSession)

		// GCP IAM UserAccount OAuth sessions
		v1.GET("gcp/iam-user-account-oauth-url", contr.GetGcpOauthUrl)
		v1.POST("gcp/iam-user-account-oauth-sessions", contr.CreateGcpIamUserAccountOauthSession)
		v1.GET("gcp/iam-user-account-oauth-sessions/:id", contr.GetGcpIamUserAccountOauthSession)
		v1.PUT("gcp/iam-user-account-oauth-sessions/:id", contr.EditGcpIamUserAccountOauthSession)
		v1.POST("gcp/iam-user-account-oauth-sessions/:id/start", contr.StartGcpIamUserAccountOauthSession)
		v1.POST("gcp/iam-user-account-oauth-sessions/:id/stop", contr.StopGcpIamUserAccountOauthSession)
		v1.DELETE("gcp/iam-user-account-oauth-sessions/:id", contr.DeleteGcpIamUserAccountOauthSession)
		v1.GET("gcp/named-configurations", contr.GetNamedConfigurations)

		// WebSocket
		v1.GET("ws", contr.GetWs)
	}
}
