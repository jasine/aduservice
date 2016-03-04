package basic

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io/ioutil"
	"log"
	"sync"
)

const (
	Init_UserNmae = "admin"
	Init_UserPwd  = "admin"
)

type BasicAdu struct {
	localmd5 []byte
	context  interface{}
	lock     sync.Mutex
	loadfunc func(interface{}) ([]byte, error)
	savefunc func(interface{}, []byte) error
}

func (b *BasicAdu) SetLoadFunc(load func(interface{}) ([]byte, error)) {
	b.loadfunc = load
}

func (b *BasicAdu) SetSaveFunc(save func(interface{}, []byte) error) {
	b.savefunc = save
}

func NewBasicAdu(context interface{}) *BasicAdu {
	b := &BasicAdu{
		context:  context,
		loadfunc: GetUserItemFromFile,
		savefunc: SetuserItemToFile,
	}
	if _, err := b.GetLocalMd5(); err != nil {
		log.Println(err.Error())
	}
	return b
}

func GetUserItemFromFile(context interface{}) ([]byte, error) {
	filename, ok := context.(string)
	if !ok {
		return nil, errors.New("filename param is not a string")
	}
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func SetuserItemToFile(context interface{}, user []byte) error {
	filename, ok := context.(string)
	if !ok {
		return errors.New("filename param is not a string")
	}
	err := ioutil.WriteFile(filename, user, 0600)
	return err
}

func (b *BasicAdu) GetLocalMd5() ([]byte, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	// bs, err := ioutil.ReadFile(b.file)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }
	bs, err := b.loadfunc(b.context)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	b.localmd5 = bs
	return bs, nil
}

func (b *BasicAdu) setLocalMd5() error {
	b.lock.Lock()
	defer b.lock.Unlock()
	// err := ioutil.WriteFile(b.file, b.localmd5, 0600)
	err := b.savefunc(b.context, b.localmd5)
	return err
}

func ComputeMd5(name, pwd string) []byte {
	buf := make([]byte, 0, 128)
	buf = append(buf, []byte(name)...)
	buf = append(buf, []byte(":")...)
	buf = append(buf, []byte(pwd)...)
	buf = append(buf, []byte("#deepglint")...)
	md5Ctx := md5.New()
	if _, err := md5Ctx.Write(buf); err != nil {
		log.Println(err.Error())
		return []byte("error") //todo:throw err
	}
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
	if pass && (len(newpwd) < 6 || len(newpwd) > 16) {
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
