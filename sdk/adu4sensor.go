package sdk

import (
	"github.com/deepglint/aduservice/basic"
)

var filePath string = "/data/adu/auth"

var basicAuth *basic.BasicAdu = basic.NewBasicAdu(filePath)

func BasicAuth(username, password string) bool {
	return basicAuth.Auth(username, password)
}

func ChangePwd(username, oldpwd, newpwd string) bool {
	return basicAuth.ChangePwd(username, oldpwd, newpwd)
}

func ResetUserAndPwd() bool {
	return basicAuth.ResetUserAndPwd()
}
