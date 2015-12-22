package main

import (
	//"encoding/hex"
	"fmt"
	"github.com/deepglint/aduservice/basic"
	"github.com/vulcand/oxy/utils"
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
	HttpAddr := ":8186"
	adu := basic.NewBasicAdu(file_name)
	cl := &Controller{
		adu: adu,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", cl.Auth)
	mux.HandleFunc("/changepwd", cl.Changepwd)
	mux.HandleFunc("/reset", cl.Reset)
	mux.HandleFunc("/update", cl.Update)

	log.Printf("Http server listens on %s\n", HttpAddr)
	http.ListenAndServe(HttpAddr, mux)
}

type Controller struct {
	adu *basic.BasicAdu
}

func (c *Controller) praseAndAuth(r *http.Request) (bool, string, string) {
	auth, err := utils.ParseAuthHeader(r.Header.Get("Authorization"))

	if err == nil && c.adu.Auth(auth.Username, auth.Password) {
		return true, auth.Username, auth.Password
	}
	return false, "", ""
}

// only for vulcand
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
	if c.adu.ChangePwd(name, pwd, newpwd) {
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
	if c.adu.ResetUserAndPwd() {
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
