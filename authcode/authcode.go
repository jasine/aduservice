package authcode

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/deepglint/aduservice/util"
	"github.com/deepglint/glog"
	"strconv"
	"strings"
	"time"
)

var (
	ServerUID = "deepglint"
)

func SetServerUid(uid string) {
	ServerUID = uid
}

func GenAuthCode() (string, error) {
	d := time.Now().UnixNano()
	data := fmt.Sprintf("code:%s:%d", ServerUID, d)
	dataEncrypt, err := util.Encrypt([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(dataEncrypt), nil
}

func AuthCodePair(codeclient, codedeepglint string) bool {
	uidclient, dclient, err := checkcode(codeclient, "code")
	if err != nil {
		glog.Error(err.Error())
		return false
	}
	uiddeepglint, ddeepglint, err := checkcode(codedeepglint, "pair")
	if err != nil {
		glog.Error(err.Error())
		return false
	}
	if uidclient != uiddeepglint || uidclient != ServerUID || ddeepglint.Sub(dclient) >= time.Second*3600 || ddeepglint.Sub(dclient) < 0 {
		return false
	}
	return true
}

func GenPairAuthCode(code string) (string, error) {
	uid, _, err := checkcode(code, "code")
	if err != nil {
		return "", err
	}
	d := time.Now().UnixNano()
	data := fmt.Sprintf("pair:%s:%d", uid, d)
	dataEncrypt, err := util.Encrypt([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(dataEncrypt), nil
}

func checkcode(code string, prefix string) (string, time.Time, error) {
	dataEncrypt, err := hex.DecodeString(code)
	if err != nil {
		return "", time.Now(), err
	}
	data, err := util.Decrypt(dataEncrypt)
	if err != nil {
		return "", time.Now(), err
	}
	ss := strings.Split(string(data), ":")
	if len(ss) != 3 || ss[0] != prefix {
		fmt.Println(len(ss) != 3, ss[0], prefix)
		return "", time.Now(), errors.New("bad code format")
	}
	uid := ss[1]
	unixtimenano, err := strconv.ParseInt(ss[2], 10, 64)
	if err != nil {
		return "", time.Now(), err
	}
	dd := time.Unix(0, unixtimenano)
	//fmt.Println(dd)
	if time.Now().Sub(dd) >= time.Second*3600 || time.Now().Sub(dd) < 0 {
		return "", time.Now(), errors.New("code timeout")
	}
	return uid, dd, nil
}
