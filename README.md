# aduservice
The ADU Service is the management service for the 
Auth, Device, Developer, Enduser

NOTE:only for go1.4, util/encryption does not work in go1.5

## 服务
### basic auth - sensor

* 需要将sdk/adu4sensor.go嵌入到bumble中编译运行，无需单独部署
* 依赖密码文件，见basic/auth目录，运行时放置到指定目录，默认路径为"config/auth"
	- InitBasicAuth4Sensor 初始化adu
	- ChangeFilePath 修改密码文件路径
	- BasicAuth 验证
	- ChangePwd 修改密码
	- ResetUserAndPwd 重置密码文件
	- GetVersion 版本号

### basic auth - server

* server端运行aduservice服务，提供用户名密码验证、修改密码和重置密码的接口，需要单独部署
* 如需要通过aduservice对server端的http接口进行验证，需要将sdk/adu4vulcand.go嵌入到vulcand中，重新编译并配置
* 详细说明见app/aduserver目录

## 部署

### aduservice server端部署
```
$ ./aduservice
依赖 /data/adu/auth
```
### vulcand server端部署
```
$ ./vulcand.linux --apiInterface=0.0.0.0 --etcd=http://127.0.0.1:4001
```
### vulcand中间件配置
```
$ ./vctl   backend upsert -id b1
$ ./vctl/vctl server upsert -id srv1 -b b1 -url http://localhost:8186
$ ./vctl/vctl frontend upsert -id f1 -b b1 -route 'PathRegexp("/.*")'
$ ./vctl/vctl basicauth4t2  upsert -f f1 -id m1
```

## 接口
```
server : http://<server_ip:port>:/api/<name>
bumble : http://<bumble_ip:port>:/api/set/<name>
```

* 登录(该接口不需要basic auth 验证)

```
Method: POST 

Name: login 

Body: name:pwd ps. admin:admin 

Response: string "true" or error msg
```

* 设置新的密码(该接口不需要basic auth 验证)

```
Method: POST 

Name: changepwd 

Body: name:pwd:newpwd ps. admin:admin:admin1 

Response: string "SUCCESS" or error msg
```

* 重置用户名和密码为admin:admin(需要basic auth 验证)

```
Method: POST 

Name: resetpwd 

Response: string "SUCCESS" or error msg
```
