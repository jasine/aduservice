package basic

import (
	"math/rand"
	"testing"
)

func TestBasic(t *testing.T) {
	adu := NewBasicAdu("/data/adu/auth")
	if _, err := adu.ResetUserAndPwd(); err != nil {
		t.Error(err.Error())
	}
	if ok, _ := adu.Auth("admin", "admin"); !ok {
		t.Error("auth false")
	}
}

func TestNoFile(t *testing.T) {
	adu := NewBasicAdu("./auth1")
	if ok, _ := adu.Auth("admin", "admin"); ok {
		t.Error("auth false")
	}
}

func TestWrongUP(t *testing.T) {
	adu := NewBasicAdu("./auth")
	if _, err := adu.ResetUserAndPwd(); err != nil {
		t.Error(err.Error())
	}
	if ok, _ := adu.Auth("admin1", "admin"); ok {
		t.Error("auth false")
	}
	if ok, _ := adu.Auth("admin", "admin1"); ok {
		t.Error("auth false")
	}
}

func TestEmptyUP(t *testing.T) {
	adu := NewBasicAdu("./auth")
	if _, err := adu.ResetUserAndPwd(); err != nil {
		t.Error(err.Error())
	}
	if ok, _ := adu.Auth("", "admin"); ok {
		t.Error("auth false")
	}
	if ok, _ := adu.Auth("admin", ""); ok {
		t.Error("auth false")
	}
	if ok, _ := adu.Auth("", ""); ok {
		t.Error("auth false")
	}
}

func TestChangPWD(t *testing.T) {
	adu := NewBasicAdu("./auth")
	if _, err := adu.ResetUserAndPwd(); err != nil {
		t.Error(err.Error())
	}
	if ok, _ := adu.Auth("admin", "admin"); !ok {
		t.Error("auth false")
	}
	if ok, _ := adu.ChangePwd("admin", "admin", "admin1"); !ok {
		t.Error("fail")
	}
	if ok, _ := adu.Auth("admin", "admin"); ok {
		t.Error("auth false")
	}
	if ok, _ := adu.Auth("admin", "admin1"); !ok {
		t.Error("auth false")
	}
}

func TestThread(t *testing.T) {
	adu := NewBasicAdu("./auth")
	if _, err := adu.ResetUserAndPwd(); err != nil {
		t.Error(err.Error())
	}
	stop := make(chan bool)
	for i := 0; i < 1000; i++ {
		go oneauth(adu, t)
		if i == 999 {
			close(stop)
		}
	}
	<-stop
}

func oneauth(adu *BasicAdu, t *testing.T) {
	if rand.Int31()%2 == 0 {
		if ok, _ := adu.Auth("admin", "admin"); !ok {
			t.Error("should be true")
		}
	} else {
		if ok, _ := adu.Auth("admin1", "admin"); ok {
			t.Error("should be false")
		}
	}
}

func readT(adu *BasicAdu, t *testing.T) {
	_, e := adu.GetLocalMd5()
	if e != nil {
		t.Error(e)
	}
}

func writeT(adu *BasicAdu, t *testing.T) {
	b, _ := adu.ResetUserAndPwd()
	if b != true {
		t.Error("bad")
	}
}

func TestReadFile(t *testing.T) {
	adu := NewBasicAdu("./auth")
	if _, err := adu.ResetUserAndPwd(); err != nil {
		t.Error(err.Error())
	}
	stop := make(chan bool)
	for i := 0; i < 1000; i++ {
		go readT(adu, t)
		if i == 999 {
			close(stop)
		}
	}
	<-stop
}

func TestWriteFile(t *testing.T) {
	adu := NewBasicAdu("./auth")
	if _, err := adu.ResetUserAndPwd(); err != nil {
		t.Error(err.Error())
	}
	stop := make(chan bool)
	for i := 0; i < 1000; i++ {
		go writeT(adu, t)
		if i == 999 {
			close(stop)
		}
	}
	<-stop
}
