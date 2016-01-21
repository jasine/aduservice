package util

import (
	//"bytes"
	"math/rand"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestBasicEncrypt(t *testing.T) {
	data := []byte("hello world")
	round(t, data)
}

func TestValid(t *testing.T) {
	data := []byte("hello world")
	data1, err := Encrypt(data)
	if err != nil {
		t.Error(err)
	}
	data1[0] = byte(128)
	_, err = Decrypt(data1)
	if err == nil {
		t.Error(err)
	}
	//equal(t, data, data2)
}

func TestBigValid(t *testing.T) {
	data := make([]byte, 1024*1024)
	data[1024] = byte(110)
	data1, err := Encrypt(data)
	if err != nil {
		t.Error(err)
	}
	data1[0] = byte(128)
	_, err = Decrypt(data1)
	if err == nil {
		t.Error(err)
	}
	//equal(t, data, data2)
}

func round(t *testing.T, data []byte) {
	data1, err := Encrypt(data)
	if err != nil {
		t.Error(err)
	}
	data2, err := Decrypt(data1)
	if err != nil {
		t.Error(err)
	}
	equal(t, data, data2)
}

func TestRound(t *testing.T) {
	data := make([]byte, 0, 1024)
	for i := 0; i < 1024*10; i++ {
		b := byte(int8(rand.Uint32() % 256))
		data = append(data, b)
		round(t, data)
	}
}

func TestBigRound(t *testing.T) {
	set := [...]int{1024 * 1024, 1024 * 1024 * 10, 1024 * 1024 * 100, 1024 * 1024 * 1024}
	for i := 0; i < len(set); i++ {
		data := make([]byte, set[i])
		data[1024] = byte(128)
		go round(t, data)
	}
}

func equal(t *testing.T, act, exp interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n",
			filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}
