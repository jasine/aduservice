# aduservice
The ADU Service is the management service for the 
Auth, Device, Developer, Enduser

## Version - 0.0.1.2

修改login接口来兼容t1, 验证失败返回错误信息和400, 验证成功返回“true”和200

## Changlog

## Version - 0.0.1.1

添加了基于 auth code 和 pair code 的重置功能

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
$ sudo docker run -i -t -d -p 127.0.0.1:8186:8186 -v /data/adu:/data/adu --name adu --restart=always 192.168.5.46:5000/aduservice-t2-test:0.0.1.2 ./aduservice.linux

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

## 测试环境
* vulcand 中间件 : http://192.168.2.31:8181/test
* aduservice : http://192.168.2.31:8186/api/login
