package basic

import (
	"math/rand"
	"testing"
)

func TestBasic(t *testing.T) {
	adu := NewBasicAdu("/data/adu/auth")
	adu.ResetUserAndPwd()
	if !adu.Auth("admin", "admin") {
		t.Error("auth false")
	}
}

func TestNoFile(t *testing.T) {
	adu := NewBasicAdu("./auth1")
	if adu.Auth("admin", "admin") {
		t.Error("auth false")
	}
}

func TestWrongUP(t *testing.T) {
	adu := NewBasicAdu("./auth")
	adu.ResetUserAndPwd()
	if adu.Auth("admin1", "admin") {
		t.Error("auth false")
	}
	if adu.Auth("admin", "admin1") {
		t.Error("auth false")
	}
}

func TestEmptyUP(t *testing.T) {
	adu := NewBasicAdu("./auth")
	adu.ResetUserAndPwd()
	if adu.Auth("", "admin") {
		t.Error("auth false")
	}
	if adu.Auth("admin", "") {
		t.Error("auth false")
	}
	if adu.Auth("", "") {
		t.Error("auth false")
	}
}

func TestChangPWD(t *testing.T) {
	adu := NewBasicAdu("./auth")
	adu.ResetUserAndPwd()
	if !adu.Auth("admin", "admin") {
		t.Error("auth false")
	}
	adu.ChangePwd("admin", "admin", "admin1")
	if adu.Auth("admin", "admin") {
		t.Error("auth false")
	}
	if !adu.Auth("admin", "admin1") {
		t.Error("auth false")
	}
}

func TestThread(t *testing.T) {
	adu := NewBasicAdu("./auth")
	adu.ResetUserAndPwd()
	stop := make(chan bool)
	for i := 0; i < 100000; i++ {
		go oneauth(adu, t)
		if i == 99999 {
			close(stop)
		}
	}
	<-stop
}

func oneauth(adu *BasicAdu, t *testing.T) {
	if rand.Int31()%2 == 0 {
		if !adu.Auth("admin", "admin") {
			t.Error("should be true")
		}
	} else {
		if adu.Auth("admin1", "admin") {
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
	b := adu.ResetUserAndPwd()
	if b != true {
		t.Error("bad")
	}
}

func TestReadFile(t *testing.T) {
	adu := NewBasicAdu("./auth")
	adu.ResetUserAndPwd()
	stop := make(chan bool)
	for i := 0; i < 100000; i++ {
		go readT(adu, t)
		if i == 99999 {
			close(stop)
		}
	}
	<-stop
}

func TestWriteFile(t *testing.T) {
	adu := NewBasicAdu("./auth")
	adu.ResetUserAndPwd()
	stop := make(chan bool)
	for i := 0; i < 100000; i++ {
		go writeT(adu, t)
		if i == 99999 {
			close(stop)
		}
	}
	<-stop
}
