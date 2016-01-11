package main

import (
	//"encoding/hex"
	"flag"
	"fmt"
	"github.com/deepglint/aduservice/basic"
	"github.com/vulcand/oxy/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	file_name  = "/data/adu/auth"
	TRUE_BODY  = "SUCCESS"
	FALSE_BODY = "FAIL"
)

func main() {
	HttpAddr := ""
	flag.StringVar(&HttpAddr, "port", ":8186", "http server port p.s. :8186")
	flag.Parse()
	adu := basic.NewBasicAdu(file_name)
	cl := &Controller{
		adu: adu,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth_basic", cl.Auth)
	mux.HandleFunc("/changepwd_basic", cl.Changepwd)
	mux.HandleFunc("/api/resetpwd", cl.Reset)
	mux.HandleFunc("/api/login", cl.LoginNoBasic)
	mux.HandleFunc("/api/changepwd", cl.ChangepwdNoBasic)
	//	mux.HandleFunc("/api/resetpwd", cl.ResetNoBasic)
	mux.HandleFunc("/update", cl.Update)
	mux.HandleFunc("/test", Test)

	log.Printf("Http server listens on %s\n", HttpAddr)
	http.ListenAndServe(HttpAddr, mux)
}
func Test(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, TRUE_BODY)
}

type Controller struct {
	adu *basic.BasicAdu
}

func (c *Controller) LoginNoBasic(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		fmt.Fprint(w, "empty body")
		return
	}
	up := strings.Split(string(body), ":")
	if len(up) != 2 {
		fmt.Fprint(w, "bad post body format")
		return
	}
	authpass, err := c.adu.Auth(up[0], up[1])
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	if !authpass {
		fmt.Fprint(w, "auth fail")
		return
	}
	fmt.Fprint(w, TRUE_BODY)
}

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
	fmt.Fprint(w, TRUE_BODY)
}

func (c *Controller) ResetNoBasic(w http.ResponseWriter, r *http.Request) {
	b, err := c.adu.ResetUserAndPwd()
	if err == nil && b {
		fmt.Fprint(w, TRUE_BODY)
		return
	}
	fmt.Fprint(w, FALSE_BODY)
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

func (c *Controller) Auth(w http.ResponseWriter, r *http.Request) {
	pass, _, _ := c.praseAndAuth(r)
	if pass {
		fmt.Fprint(w, TRUE_BODY)
	} else {
		fmt.Fprint(w, FALSE_BODY)
	}
}

func (c *Controller) Changepwd(w http.ResponseWriter, r *http.Request) {
	pass, name, pwd := c.praseAndAuth(r)
	if !pass {
		fmt.Fprint(w, FALSE_BODY)
		return
	}
	r.ParseForm()
	var newpwd string
	if v, ok := r.Form["pwd"]; ok {
		newpwd = strings.Join(v, "")
	} else {
		fmt.Fprint(w, FALSE_BODY)
		return
	}
	b, err := c.adu.ChangePwd(name, pwd, newpwd)
	if b && err == nil {
		fmt.Fprint(w, TRUE_BODY)
		return
	}
	fmt.Fprint(w, FALSE_BODY)
}

func (c *Controller) Reset(w http.ResponseWriter, r *http.Request) {
	pass, _, _ := c.praseAndAuth(r)
	if !pass {
		fmt.Fprint(w, FALSE_BODY)
		return
	}
	b, err := c.adu.ResetUserAndPwd()
	if err == nil && b {
		fmt.Fprint(w, TRUE_BODY)
		return
	}
	fmt.Fprint(w, FALSE_BODY)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	bs, err := c.adu.GetLocalMd5()
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, FALSE_BODY)
		return
	}
	fmt.Fprint(w, fmt.Sprintf("%s", string(bs)))
}
