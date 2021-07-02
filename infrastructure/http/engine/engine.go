package engine

import (
	"fmt"
	"leapp_daemon/infrastructure/http/middleware"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller"
	"leapp_daemon/providers"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	contr := controller.EngineController{Providers: providers}

	v1 := ginEngine.Group("/api/v1")
	{
		// All sessions
		v1.GET("sessions", contr.ListSession)

		// AWS sessions
		v1.GET("aws/named-profiles", contr.ListNamedProfiles)
		v1.GET("aws/regions", contr.GetAwsRegionList)
		v1.PUT("aws/sessions/:id/region", contr.EditAwsRegion)

		// AWS IAM User sessions
		v1.GET("aws/iam-user-sessions/:id", contr.GetAwsIamUserSession)
		v1.POST("aws/iam-user-sessions", contr.CreateAwsIamUserSession)
		v1.PUT("aws/iam-user-sessions/:id", contr.UpdateAwsIamUserSession)
		v1.DELETE("aws/iam-user-sessions/:id", contr.DeleteAwsIamUserSession)
		v1.POST("aws/iam-user-sessions/:id/confirm-mfa-token", contr.ConfirmMfaToken)
		v1.POST("aws/iam-user-sessions/:id/start", contr.StartAwsIamUserSession)
		v1.POST("aws/iam-user-sessions/:id/stop", contr.StopAwsIamUserSession)

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

		// Alibaba RAM User sessions
		v1.GET("/alibaba/ram-user-sessions/:id", contr.GetAlibabaRamUserSessionController)
		v1.POST("/alibaba/ram-user-sessions", contr.CreateAlibabaRamUserSessionController)
		v1.PUT("/alibaba/ram-user-sessions/:id", contr.UpdateAlibabaRamUserSessionController)
		v1.DELETE("/alibaba/ram-user-sessions/:id", contr.DeleteAlibabaRamUserSessionController)
		v1.POST("/alibaba/ram-user-sessions/:id/start", contr.StartAlibabaRamUserSessionController)
		v1.POST("/alibaba/ram-user-sessions/:id/stop", contr.StopAlibabaRamUserSessionController)

		// Alibaba RAM Role Federated sessions
		v1.GET("/alibaba/ram-role-federated-sessions/:id", contr.GetAlibabaRamRoleFederatedSessionController)
		v1.POST("/alibaba/ram-role-federated-sessions", contr.CreateAlibabaRamRoleFederatedSessionController)
		v1.PUT("/alibaba/ram-role-federated-sessions/:id", contr.EditAlibabaRamRoleFederatedSessionController)
		v1.DELETE("/alibaba/ram-role-federated-sessions/:id", contr.DeleteAlibabaRamRoleFederatedSessionController)
		v1.POST("/alibaba/ram-role-federated-sessions/:id/start", contr.StartAlibabaRamRoleFederatedSessionController)
		v1.POST("/alibaba/ram-role-federated-sessions/:id/stop", contr.StopAlibabaRamRoleFederatedSessionController)

		// Alibaba RAM Role Chained sessions
		v1.GET("/alibaba/ram-role-chained-sessions/:id", contr.GetAlibabaRamRoleChainedSessionController)
		v1.POST("/alibaba/ram-role-chained-sessions", contr.CreateAlibabaRamRoleChainedSessionController)
		v1.PUT("/alibaba/ram-role-chained-sessions/:id", contr.EditAlibabaRamRoleChainedSessionController)
		v1.DELETE("/alibaba/ram-role-chained-sessions/:id", contr.DeleteAlibabaRamRoleChainedSessionController)
		v1.POST("/alibaba/ram-role-chained-sessions/:id/start", contr.StartAlibabaRamRoleChainedSessionController)
		v1.POST("/alibaba/ram-role-chained-sessions/:id/stop", contr.StopAlibabaRamRoleChainedSessionController)

		// WebSocket
		v1.GET("ws", contr.GetWs)
	}
}
