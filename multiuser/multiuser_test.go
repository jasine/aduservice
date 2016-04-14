package multiuser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMultiUser(t *testing.T) {
	m := NewMultiUserController()
	Convey("测试文件存储", t, func() {

		Convey("打开测试账号文件test.auth", func() {
			err := m.LoadFileToMap("test.auth")
			So(err, ShouldBeNil)
			Convey("添加账号admin:admin", func() {
				err = m.AddUser("admin", "admin", []string{})
				So(err, ShouldNotBeNil)
				Convey("将账号存储在本地", func() {
					err = m.SaveToFile("test.auth")
					So(err, ShouldBeNil)
					Convey("重新load文件", func() {
						err = m.LoadFileToMap("test.auth")
						So(err, ShouldBeNil)
						Convey("测试auth是否可以登录", func() {
							err = m.AddUser("admin", "admin", []string{})
							So(err, ShouldNotBeNil)
						})
						// err = m.AddUser("admin", "admin", []string{})
						// ShouldBeNil(err)
					})

				})

			})

		})

	})

	Convey("测试新增用户", t, func() {

		Convey("打开测试账号文件test.auth", func() {
			So(m.LoadFileToMap("test.auth"), ShouldBeNil)
			Convey("添加账号admin:admin", func() {
				So(m.AuthUser("admin", "admin"), ShouldNotBeNil)
				So(m.AddUser("admin", "admin", []string{}), ShouldNotBeNil)
				Convey("添加测试账号admin1:admin1", func() {
					So(m.AddUser("admin1", "admin1", []string{}), ShouldNotBeNil)
					Convey("使用admin:admin登录", func() {

						Convey("使用admin1:admin1登录", func() {
							So(m.AuthUser("admin1", "admin1"), ShouldBeNil)
							Convey("保存file并退出", nil)
							So(m.SaveToFile("test.auth"), ShouldBeNil)
						})

					})

				})

			})

		})

	})

	Convey("测试修改Auth", t, func() {
		err := m.LoadFileToMap("test.auth")
		So(err, ShouldBeNil)
		So(m.UpdateUserAuth("admin", "admin1"), ShouldBeNil)
		So(m.AuthUser("admin", "admin1"), ShouldBeNil)
		So(m.AuthUser("admin", "admin"), ShouldNotBeNil)
		So(m.SaveToFile("test.auth"), ShouldBeNil)
		So(m.LoadFileToMap("test.auth"), ShouldBeNil)
		So(m.AuthUser("admin", "admin1"), ShouldBeNil)
		So(m.AuthUser("admin", "admin"), ShouldNotBeNil)
	})

	Convey("测试修改权限", t, nil)

}
