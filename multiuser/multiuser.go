package multiuser

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"errors"
	"log"
	"os"
	"sync"
)

// const (
// 	Init_UserName = "admin"
// 	Init_UserPwd  = "admin"
// )

type MultiUserController struct {
	lock  sync.Mutex
	users map[string]*User
	//loadfunc func(interface{}) ([]byte, error)
	//savefunc func(interface{}, []byte) error
}

type User struct {
	Md5      []byte
	Username string
	Right    []string
}

func NewMultiUserController() *MultiUserController {
	m := new(MultiUserController)
	m.users = make(map[string]*User)
	return m
}

func (this *MultiUserController) ChangeRight(name string, right []string) error {
	if this.users[name] == nil {
		return errors.New("User not Existed")
	}
	this.lock.Lock()
	this.users[name].Right = right
	this.lock.Unlock()
	return nil
}

func (this *MultiUserController) AuthUser(name, pwd string) error {
	if this.users[name] == nil {
		return errors.New("User not Existed")
	}
	md := ComputeMd5(name, pwd)
	if !bytes.Equal(this.users[name].Md5, md) {
		return errors.New("Auth Failed")
	}
	return nil
}

func (this *MultiUserController) AddUser(name, pwd string, right []string) error {
	if this.users[name] != nil {
		return errors.New("User Existed")
	}
	md := ComputeMd5(name, pwd)
	this.lock.Lock()
	this.users[name] = &User{Md5: md, Username: name, Right: right}
	this.lock.Unlock()
	return nil
}

func (this *MultiUserController) UpdateUserAuth(name, pwd string) error {
	if this.users[name] == nil {
		return errors.New("User not Existed")
	}
	md := ComputeMd5(name, pwd)
	this.lock.Lock()
	this.users[name].Md5 = md
	this.lock.Unlock()
	return nil
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

func (this *MultiUserController) SaveToFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(this.users)
	return err
}

func (this *MultiUserController) LoadFileToMap(file string) error {
	//m := make(map[string]interface{}, 0)
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewDecoder(f)
	err = enc.Decode(&this.users)
	if err != nil {
		return err
	}
	return nil
}
