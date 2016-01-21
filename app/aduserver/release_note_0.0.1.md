# aduservice
The ADU Service is the management service for the 
Auth, Device, Developer, Enduser

## Version - 0.0.1.1

添加了基于 auth code 和 pair code 的重置功能

## Changlog

### version 0.0.1

实现了basic auth 验证服务的API和SDK

## 服务
### basic auth sdk for sensor

* 需要将sdk/adu4sensor.go嵌入到bumble中编译运行，无需单独部署
* 依赖 config/auth 文件

### basic auth - server

* server端运行aduservice服务，提供用户名密码验证、修改密码和重置密码的接口，需要单独部署
* 依赖 /data/adu/auth 文件

### basic auth sdk for vulcand

* 如需要通过aduservice对server端的http接口进行验证，需要将sdk/adu4vulcand.go嵌入到vulcand中，重新编译并配置
* 依赖 basic auth - server
* 依赖 etcd

## 部署

### aduservice server端部署
```
$ sudo docker run -i -t -d -p 8186:8186 -v /data/adu:/data/adu --name adu 192.168.5.46:5000/aduservice-t2-test:0.0.1.1 ./aduservice.linux

```
### vulcand server端部署
```
etcd
$ sudo docker run --net=host --restart=always -i -t -d --name etcd  -v /data/etcd:/tmp  -v /etc/localtime:/etc/localtime:ro 192.168.5.46:5000/etcd:2.1.1  --data-dir /tmp/dg --name dg  -advertise-client-urls http://127.0.0.1:4001 -listen-client-urls http://127.0.0.1:4001 -initial-advertise-peer-urls http://0.0.0.0:2380 -listen-peer-urls http://0.0.0.0:2380 -initial-cluster-token etcd-cluster -initial-cluster dg=http://0.0.0.0:2380 -initial-cluster-state new


vulcand
$ sudo docker run -t -d -i --restart=always --net=host --name=vulcand 192.168.5.46:5000/vulcand-t2-test:0.8.12 ./vulcand --apiInterface=0.0.0.0 --etcd=http://127.0.0.1:4001

```
### vulcand中间件配置
```
$ ./vctl   backend upsert -id b1
$ ./vctl server upsert -id srv1 -b b1 -url http://localhost:8186
$ ./vctl frontend upsert -id f1 -b b1 -route 'PathRegexp("/.*")'
$ ./vctl basicauth4t2  upsert -f f1 -id m1
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

Response: string "SUCCESS" or error msg
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

* 获取authcode(该接口不需要basic auth 验证，只在server端提供)

```
Method: Get 

Name: authcode 

Body: 

Response: string authcode or error msg
```

* 获取paircode(该接口不需要basic auth 验证，只在server端提供，不提供给客户使用，客户联系我们后，我们提供)

```
Method: Get 

Param: code 例如：http://localhost:8186/api/paircode?code=140c63838944357c7a6faff182ac167bbfb0966c5f9177e60def74a88941d3da2b442487a71656faba8766445e7564101886f57ef2092934ea27ce7fb8f66cbd

Name: paircode 

Body: 

Response: string paircode or error msg
```

* 重置密码(该接口不需要basic auth 验证，只在server端提供，供忘记密码，重置密码使用)

```
Method: Get 
Param: code pair 例如：http://localhost:8186/api/resetpwd_code?code=140c63838944357c7a6faff182ac167bbfb0966c5f9177e60def74a88941d3da2b442487a71656faba8766445e7564101886f57ef2092934ea27ce7fb8f66cbd&pair=f28298040d3a17e499cd7ceee4952618482501cee4df381f20d1f974006061807ac877128539509aaf98cfe2eb541b489290f573052780a9806e0f9d2a5bf09b

Name: resetpwd_code 

Body: 

Response: string "SUCCESS" or error msg
```


## 测试环境
* vulcand 中间件 : http://192.168.2.26:8181/test
* aduservice : http://192.168.2.26:8186/api/login
