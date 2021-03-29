package g_suite_auth_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type GSuiteAuthSecondStepRequestDto struct {
	Captcha string `json:"captcha"`
	CaptchaInputId string `json:"captchaInputId"`
	CaptchaUrl string `json:"captchaUrl"`
	CaptchaForm string `json:"captchaForm"`
	Password string `json:"password"`
	LoginForm string `json:"loginForm"`
	LoginUrl string `json:"loginUrl"`
}

func (requestDto *GSuiteAuthSecondStepRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}