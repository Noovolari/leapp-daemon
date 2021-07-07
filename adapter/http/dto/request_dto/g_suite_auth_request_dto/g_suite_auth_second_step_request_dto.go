package g_suite_auth_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

type GSuiteAuthSecondStepRequestDto struct {
	Captcha        string `json:"captcha"`
	CaptchaInputId string `json:"captchaInputId"`
	CaptchaUrl     string `json:"captchaUrl"`
	CaptchaForm    string `json:"captchaForm"`
	Password       string `json:"password"`
	LoginForm      string `json:"loginForm"`
	LoginUrl       string `json:"loginUrl"`
}

func (requestDto *GSuiteAuthSecondStepRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
