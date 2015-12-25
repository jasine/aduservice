package sdk

import (
	"encoding/base64"
	"github.com/deepglint/aduservice/basic"
	"github.com/go-martini/martini"
	"net/http"
	"strings"
)

var filePath string = "config/auth"

// User is the authenticated username that was extracted from the request.
type User string

// BasicRealm is used when setting the WWW-Authenticate response header.
var BasicRealm = "Authorization Required"

var basicAuth *basic.BasicAdu = nil

func InitBasicAuth4Sensor() {
	basicAuth = basic.NewBasicAdu(filePath)
}

func BasicAuth(username, password string) (bool, error) {
	return basicAuth.Auth(username, password)
}

func ChangePwd(username, oldpwd, newpwd string) (bool, error) {
	return basicAuth.ChangePwd(username, oldpwd, newpwd)
}

func ResetUserAndPwd() (bool, error) {
	return basicAuth.ResetUserAndPwd()
}

// BasicFunc returns a Handler that authenticates via Basic Auth using the provided function.
// The function should return true for a valid username/password combination.
func BasicFunc(authfn func(string, string) (bool, error), router []string) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		for i := 0; i < len(router); i++ {
			if req.URL.Path == router[i] {
				return
			}
		}
		auth := req.Header.Get("Authorization")
		if len(auth) < 6 || auth[:6] != "Basic " {
			unauthorized(res)
			return
		}
		b, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			unauthorized(res)
			return
		}
		tokens := strings.SplitN(string(b), ":", 2)
		if len(tokens) != 2 {
			unauthorized(res)
			return
		}
		authpass, err := authfn(tokens[0], tokens[1])
		if !authpass || err != nil {
			unauthorized(res)
			return
		}
		c.Map(User(tokens[0]))
	}
}

func unauthorized(res http.ResponseWriter) {
	res.Header().Set("WWW-Authenticate", "Basic realm=\""+BasicRealm+"\"")
	http.Error(res, "Not Authorized", http.StatusUnauthorized)
}
