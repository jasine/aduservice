package sdk

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/deepglint/aduservice/basic"
	"github.com/vulcand/oxy/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	TRUE_BODY  = "SUCCESS"
	FALSE_BODY = "FAIL"
)

type Adu4vulcand struct {
	localmd5 []byte
	stopchan chan bool
	addr     string
}

func NewAdu4vulcand() *Adu4vulcand {
	av := &Adu4vulcand{
		localmd5: nil,
		stopchan: make(chan bool),
		addr:     "127.0.0.1:8186",
	}
	return av
}

func (a *Adu4vulcand) update() error {
	u, _ := url.Parse(fmt.Sprintf("http://%s/update", a.addr))
	res, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
		return err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	res.Body.Close()
	if string(result) == FALSE_BODY {
		log.Println(FALSE_BODY)
		return errors.New("false body")
	} else {
		a.localmd5 = result
	}
	return nil
}

func (a *Adu4vulcand) Check(r *http.Request) bool {
	auth, err := utils.ParseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		log.Println(err)
		return false
	}
	if len(a.localmd5) == 0 {
		log.Println("len(a.localmd5)=0")
		return false
	}
	return bytes.Equal(basic.ComputeMd5(auth.Username, auth.Password), a.localmd5)

}

func (a *Adu4vulcand) run() {
	ticker := time.Tick(time.Duration(997) * time.Millisecond)
	for {
		select {
		case <-ticker:
			a.update()
			//log.Println(a.localmd5)
		case <-a.stopchan:
			return
		}
	}
}

func (a *Adu4vulcand) Start() {
	go a.run()
}
