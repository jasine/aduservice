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
	path  string
	//loadfunc func(interface{}) ([]byte, error)
	//savefunc func(interface{}, []byte) error
}

type User struct {
	Md5      []byte
	Username string
	Right    []string
}

func NewMultiUserController(filename string) *MultiUserController {
	m := new(MultiUserController)
	m.users = make(map[string]*User)
	m.path = filename
	err := m.LoadFileToMap()
	if err != nil {
		println("Create init super admin account admin")
		right := []string{"SuperAdmin"}
		m.AddUser("admin", "admin", right)
	}
	return m
}

func (this *MultiUserController) GetUserRight(name string) ([]string, error) {
	if name == "" {
		return []string{}, errors.New("用户名为空")
	}
	if this.users[name] == nil {
		return []string{}, errors.New("用户不存在")
	}
	// if _, err := this.AuthUser(name, pass); err != nil {
	// 	return []string{}, err
	// }
	return this.users[name].Right, nil
}

func (this *MultiUserController) ChangeRight(name string, right []string) error {
	if this.users[name] == nil {
		return errors.New("User not Existed")
	}
	this.lock.Lock()
	this.users[name].Right = right
	this.SaveToFile()
	this.lock.Unlock()
	return nil
}

func (this *MultiUserController) GetAllUsers() map[string]*User {
	return this.users
}

func (this *MultiUserController) AuthUser(name, pwd string) (bool, error) {
	if this.users[name] == nil {
		return false, errors.New("User not Existed")
	}
	md := ComputeMd5(name, pwd)
	if !bytes.Equal(this.users[name].Md5, md) {
		return false, errors.New("Auth Failed")
	}
	return true, nil
}

func (this *MultiUserController) AddUser(name, pwd string, right []string) error {
	if this.users[name] != nil {
		return errors.New("User Existed")
	}
	md := ComputeMd5(name, pwd)
	this.lock.Lock()
	this.users[name] = &User{Md5: md, Username: name, Right: right}
	this.SaveToFile()
	this.lock.Unlock()
	return nil
}

func (this *MultiUserController) UpdateUserAuth(name, oldpwd, pwd string) error {
	if this.users[name] == nil {
		return errors.New("User not Existed")
	}
	if res, _ := this.AuthUser(name, oldpwd); !res {
		return errors.New("Error for Old User Auth")
	}
	md := ComputeMd5(name, pwd)
	this.lock.Lock()
	this.users[name].Md5 = md
	this.SaveToFile()
	this.lock.Unlock()
	return nil
}
func (this *MultiUserController) DeleteUser(name string) error {

	if this.users[name] == nil {
		return errors.New("User not Existed")
	}
	this.lock.Lock()
	delete(this.users, name)
	this.SaveToFile()
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

func (this *MultiUserController) SaveToFile() error {
	println(this.path)
	f, err2 := os.Create(this.path)
	if err2 != nil {
		println(err2.Error)
		return err2
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err3 := enc.Encode(this.users)
	return err3
}

func (this *MultiUserController) LoadFileToMap() error {
	//m := make(map[string]interface{}, 0)
	f, err := os.Open(this.path)
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
