package basic

import (
	"bytes"
	"crypto/md5"
	"hash"
	"io/ioutil"
)

const (
	Init_UserNmae = "admin"
	Init_UserPwd  = "admin"
)

type BasicAdu struct {
	localmd5 []byte
	file     string
	md5Ctx   hash.Hash
	buffer   *bytes.Buffer
}

func NewBasicAdu(file string) *BasicAdu {
	buf := make([]byte, 128)
	b := &BasicAdu{
		file:   file,
		md5Ctx: md5.New(),
		buffer: bytes.NewBuffer(buf),
	}
	b.getLocalMd5()
	return b
}

func (b *BasicAdu) getLocalMd5() error {
	bs, err := ioutil.ReadFile(b.file)
	if err != nil {
		return err
	}
	b.localmd5 = bs
	return nil
}

func (b *BasicAdu) setLocalMd5() error {
	err := ioutil.WriteFile(b.file, b.localmd5, 0600)
	return err
}

func (b *BasicAdu) computeMd5(name, pwd string) []byte {
	b.buffer.Reset()
	b.buffer.WriteString(name)
	b.buffer.WriteString(":")
	b.buffer.WriteString(pwd)
	b.buffer.WriteString("#deepglint")
	b.md5Ctx.Write(b.buffer.Bytes())
	md5bs := b.md5Ctx.Sum(nil)
	return md5bs
}

func (b *BasicAdu) Auth(name, pwd string) bool {
	if b.localmd5 != nil {
		return bytes.Equal(b.localmd5, b.computeMd5(name, pwd))
	}
	return false
}

func (b *BasicAdu) ChangePwd(name, oldpwd, newpwd string) bool {
	if b.Auth(name, oldpwd) {
		b.localmd5 = b.computeMd5(name, newpwd)
		err := b.setLocalMd5()
		if err != nil {
			return false
		} else {
			return true
		}
	}
	return false
}

func (b *BasicAdu) ResetUserAndPwd() bool {
	b.localmd5 = b.computeMd5(Init_UserNmae, Init_UserPwd)
	err := b.setLocalMd5()
	if err != nil {
		return false
	} else {
		return true
	}
}
