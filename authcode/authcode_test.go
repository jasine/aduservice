package authcode

import (
	//	"math/rand"
	"testing"
)

func TestBasic(t *testing.T) {
	ServerUID = "123456"
	code, err := GenAuthCode()
	t.Log(code)
	if err != nil {
		t.Error(err.Error())
	}
	pair, err := GenPairAuthCode(code)
	if err != nil {
		t.Error(err.Error())
	}
	if !AuthCodePair(code, pair) {
		t.Error("someting error")
	}
}

func TestTwoServer(t *testing.T) {
	ServerUID = "123456"
	code, err := GenAuthCode()
	if err != nil {
		t.Error(err.Error())
	}
	ServerUID = "1234567"

	pair, err := GenPairAuthCode(code)
	if err != nil {
		t.Error(err.Error())
	}
	if AuthCodePair(code, pair) {
		t.Error("someting error")
	}
}
