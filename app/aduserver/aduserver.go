package main

import (
	//"encoding/hex"
	"flag"
	"fmt"
	"github.com/deepglint/aduservice/authcode"
	"github.com/deepglint/aduservice/basic"
	"github.com/deepglint/aduservice/controllers"
	"github.com/deepglint/muses/eventserver/models"
	"github.com/vulcand/oxy/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	fileName = "/data/adu/auth" //auth file path
	// TrueBody HTTP status
	TrueBody = "SUCCESS"
	// FalseBody HTTP bad status
	FalseBody = "FAIL"
)

func main() {
	HTTPAddr := ""
	flag.StringVar(&HTTPAddr, "port", ":8186", "http server port p.s. :8186")
	flag.Parse()
	adu := basic.NewBasicAdu(fileName)
	cl := &Controller{
		adu:           adu,
		apiController: controllers.NewApiController("api"),
	}
	cl.SetDocuments()

	//mux := http.NewServeMux()
	http.HandleFunc("/api/document", cl.apiController.Document)

	http.HandleFunc("/auth_basic", cl.Auth)
	http.HandleFunc("/changepwd_basic", cl.Changepwd)
	http.HandleFunc("/api/resetpwd", cl.Reset)
	http.HandleFunc("/api/login", cl.LoginNoBasic)
	http.HandleFunc("/api/changepwd", cl.ChangepwdNoBasic)

	http.HandleFunc("/api/authcode", cl.AuthCode)
	http.HandleFunc("/api/paircode", cl.PairCode)
	http.HandleFunc("/api/resetpwd_code", cl.ResetpwdCode)

	//	mux.HandleFunc("/api/resetpwd", cl.ResetNoBasic)
	http.HandleFunc("/update", cl.Update)
	http.HandleFunc("/test", myTest)

	log.Printf("Http server listens on %s\n", HTTPAddr)
	if err := http.ListenAndServe(HTTPAddr, nil); err != nil {
		log.Println("http serve error - ", err.Error())
	}
}

// debug for basic road
func myTest(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, TrueBody)
}

// Controller : http handler
type Controller struct {
	adu           *basic.BasicAdu
	apiController *controllers.ApiController //*ApiController
}

// AuthCode used for gen authcode-client
func (c *Controller) AuthCode(w http.ResponseWriter, r *http.Request) {
	d, e := authcode.GenAuthCode()
	if e != nil {
		fmt.Fprint(w, "error-"+e.Error())
		return
	}
	fmt.Fprint(w, d)
}

// PairCode used for gen authcode-server
func (c *Controller) PairCode(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "bad request")
		return
	}
	var code string
	if v, ok := r.Form["code"]; ok {
		code = strings.Join(v, "")
	} else {
		fmt.Fprint(w, "bad code")
		return
	}

	d, e := authcode.GenPairAuthCode(code)
	if e != nil {
		fmt.Fprint(w, "error-"+e.Error())
		return
	}
	fmt.Fprint(w, d)
}

// ResetpwdCode use auth code and pair code to reset , instead of user-pwd
func (c *Controller) ResetpwdCode(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "bad request")
		return
	}
	var code, pair string

	if v, ok := r.Form["code"]; ok {
		code = strings.Join(v, "")
	} else {
		fmt.Fprint(w, "bad code")
		return
	}

	if v, ok := r.Form["pair"]; ok {
		pair = strings.Join(v, "")
	} else {
		fmt.Fprint(w, "bad pair")
		return
	}

	pass := authcode.AuthCodePair(code, pair)
	if !pass {
		fmt.Fprint(w, "code auth fail")
		return
	}
	pass, _ = c.adu.ResetUserAndPwd()
	if !pass {
		fmt.Fprint(w, "reset fail")
		return
	}
	fmt.Fprint(w, TrueBody)
}

// LoginNoBasic Login for backbone
func (c *Controller) LoginNoBasic(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "empty body")
		return
	}
	up := strings.Split(string(body), ":")
	if len(up) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad post body format")
		return
	}
	authpass, err := c.adu.Auth(up[0], up[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if !authpass {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "auth fail")
		return
	}
	fmt.Fprint(w, "true")
}

// ChangepwdNoBasic change password for backbone
func (c *Controller) ChangepwdNoBasic(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		fmt.Fprint(w, "empty body")
		return
	}
	up := strings.Split(string(body), ":")
	if len(up) != 3 {
		fmt.Fprint(w, "bad post body format")
		return
	}
	authpass, err := c.adu.ChangePwd(up[0], up[1], up[2])
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	if !authpass {
		fmt.Fprint(w, "auth fail")
		return
	}
	fmt.Fprint(w, TrueBody)
}

// ResetNoBasic reset user-pwd for backbone
func (c *Controller) ResetNoBasic(w http.ResponseWriter, r *http.Request) {
	b, err := c.adu.ResetUserAndPwd()
	if err == nil && b {
		fmt.Fprint(w, TrueBody)
		return
	}
	fmt.Fprint(w, FalseBody)
}

func (c *Controller) praseAndAuth(r *http.Request) (bool, string, string) {
	auth, err := utils.ParseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		return false, "", ""
	}
	pass, err := c.adu.Auth(auth.Username, auth.Password)
	if err == nil && pass {
		return true, auth.Username, auth.Password
	}
	return false, "", ""
}

// Auth deprecated
func (c *Controller) Auth(w http.ResponseWriter, r *http.Request) {
	pass, _, _ := c.praseAndAuth(r)
	if pass {
		fmt.Fprint(w, TrueBody)
	} else {
		fmt.Fprint(w, FalseBody)
	}
}

// Changepwd deprecated
func (c *Controller) Changepwd(w http.ResponseWriter, r *http.Request) {
	pass, name, pwd := c.praseAndAuth(r)
	if !pass {
		fmt.Fprint(w, FalseBody)
		return
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "bad request")
		return
	}
	var newpwd string
	if v, ok := r.Form["pwd"]; ok {
		newpwd = strings.Join(v, "")
	} else {
		fmt.Fprint(w, FalseBody)
		return
	}
	b, err := c.adu.ChangePwd(name, pwd, newpwd)
	if b && err == nil {
		fmt.Fprint(w, TrueBody)
		return
	}
	fmt.Fprint(w, FalseBody)
}

// Reset reset user-pwd , basic auth
func (c *Controller) Reset(w http.ResponseWriter, r *http.Request) {
	pass, _, _ := c.praseAndAuth(r)
	if !pass {
		fmt.Fprint(w, FalseBody)
		return
	}
	b, err := c.adu.ResetUserAndPwd()
	if err == nil && b {
		fmt.Fprint(w, TrueBody)
		return
	}
	fmt.Fprint(w, FalseBody)
}

// Update used by vulcand , client should not see it
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	bs, err := c.adu.GetLocalMd5()
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, FalseBody)
		return
	}
	fmt.Fprint(w, fmt.Sprintf("%s", string(bs)))
}

func (this *Controller) SetDocuments() {

	//-----API login
	loginApi := models.NewApi("login", "/api/login", "adu", "LoginNoBasic", models.REQUEST_POST, false)

	loginApi.SetDescription(
		"login", "登录，传入 name:password，成功返回 true，失败返回错误信息")
	loginApi.SetRequest(models.BODY_TYPE_TEXT)
	//getVersionApi.AddParameter(name, paratype, demovalue, text_en, text_cn)
	loginApi.RequestBodyDemo = "admin:admin"
	//getVersionApi.SetRequest(models.BODY_TYPE_NONE)

	loginApi.SetResponse(models.BODY_TYPE_TEXT)

	/*getVersionApi.AddResponseBodyParameter("Name", models.TYPE_STRING, "getdocument",
		"the name of the API",
		"API名称",
	)*/

	loginApi.ResponseBodyDemo = "true"

	this.apiController.AppendDocument(loginApi)

	/*---------------*/

	//-----API pwdApi
	pwdApi := models.NewApi("changepwd", "/api/changepwd", "adu", "Changepwd", models.REQUEST_POST, false)

	pwdApi.SetDescription(
		"change password", "修改密码，传入 name:password:new_password，成功返回 SUCCESS，失败返回错误信息")
	pwdApi.SetRequest(models.BODY_TYPE_TEXT)
	//getVersionApi.AddParameter(name, paratype, demovalue, text_en, text_cn)
	pwdApi.RequestBodyDemo = "admin:admin:admin_pwd_new"
	//getVersionApi.SetRequest(models.BODY_TYPE_NONE)

	pwdApi.SetResponse(models.BODY_TYPE_TEXT)

	/*getVersionApi.AddResponseBodyParameter("Name", models.TYPE_STRING, "getdocument",
		"the name of the API",
		"API名称",
	)*/

	pwdApi.ResponseBodyDemo = "SUCCESS"

	this.apiController.AppendDocument(pwdApi)

	/*---------------*/

	//-----API resetApi
	resetApi := models.NewApi("resetpwd", "/api/resetpwd", "adu", "resetpwd", models.REQUEST_POST, false)

	resetApi.SetDescription(
		"reset username password", "重置用户名和密码为admin:admin(需要basic auth 验证)")
	resetApi.SetRequest(models.BODY_TYPE_TEXT)
	//getVersionApi.AddParameter(name, paratype, demovalue, text_en, text_cn)
	//resetApi.RequestBodyDemo = "admin:admin:admin_pwd_new"
	resetApi.SetRequest(models.BODY_TYPE_NONE)

	resetApi.SetResponse(models.BODY_TYPE_TEXT)

	/*getVersionApi.AddResponseBodyParameter("Name", models.TYPE_STRING, "getdocument",
		"the name of the API",
		"API名称",
	)*/

	resetApi.ResponseBodyDemo = "SUCCESS"

	this.apiController.AppendDocument(resetApi)

	/*---------------*/

}
