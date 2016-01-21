package basic

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io/ioutil"
	"log"
)

const (
	Init_UserNmae = "admin"
	Init_UserPwd  = "admin"
)

type BasicAdu struct {
	localmd5 []byte
	file     string
}

func test() {

}

func NewBasicAdu(file string) *BasicAdu {
	b := &BasicAdu{
		file: file,
	}
	b.GetLocalMd5()
	return b
}

func (b *BasicAdu) GetLocalMd5() ([]byte, error) {
	bs, err := ioutil.ReadFile(b.file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	b.localmd5 = bs
	return bs, nil
}

func (b *BasicAdu) setLocalMd5() error {
	err := ioutil.WriteFile(b.file, b.localmd5, 0600)
	return err
}

func ComputeMd5(name, pwd string) []byte {
	buf := make([]byte, 0, 128)
	buf = append(buf, []byte(name)...)
	buf = append(buf, []byte(":")...)
	buf = append(buf, []byte(pwd)...)
	buf = append(buf, []byte("#deepglint")...)
	md5Ctx := md5.New()
	md5Ctx.Write(buf)
	md5bs := md5Ctx.Sum(nil)
	return md5bs
}

func (b *BasicAdu) Auth(name, pwd string) (bool, error) {
	if b.localmd5 != nil {
		return bytes.Equal(b.localmd5, ComputeMd5(name, pwd)), nil
	}
	return false, nil
}

func (b *BasicAdu) ChangePwd(name, oldpwd, newpwd string) (bool, error) {
	pass, err := b.Auth(name, oldpwd)
	if err != nil {
		return false, err
	}
	if pass && (len(newpwd) < 1 || len(newpwd) > 128) {
		return false, errors.New("bad pwd length")
	}
	if pass {
		b.localmd5 = ComputeMd5(name, newpwd)
		err := b.setLocalMd5()
		if err != nil {
			log.Println(err)
			return false, err
		} else {
			return true, nil
		}
	}
	return false, errors.New("auth fail")
}

func (b *BasicAdu) ResetUserAndPwd() (bool, error) {
	b.localmd5 = ComputeMd5(Init_UserNmae, Init_UserPwd)
	err := b.setLocalMd5()
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		return true, nil
	}
}
