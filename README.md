### Required
- ubuntu 18.04
- go 1.13.4
- docker 19.03
- docker-compose 1.26.0
- mysql 8.0
- hyperledger fabric 1.4.6

###Ready
- Create a **github.com/fabric-app database** and import SQL
- Bring up hyperledger fabric local test net or cluster

### Config

You should modify `conf/app.ini`
Crypto-config files should be replace by you fabric network msp files
Change the network connection config file: conn-fn1.yaml, default yaml file is simple network connection file. 

```
[database]
Type = mysql
User = root
Password =
Host = 202.193.60.215:3306
Name = BC
TablePrefix = 
```
### Deploy simple network
```shell script
cd ./test/simple-network/
./network.sh  up
```
### Install
```
$ go get -u github.com/swaggo/swag/cmd/swag
$ git clone 
$ go mod tidy
$ go mod vendor
$ swag init
$ go run main.go
```
### Project information and existing API
```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (4 handlers)
[GIN-debug] POST   /api/v1/user/login        --> github.com/fabric-app/controller/api/v1.Auth (4 handlers)
[GIN-debug] POST   /api/v1/user/register     --> github.com/fabric-app/controller/api/v1.Reg (4 handlers)
[GIN-debug] GET    /api/v1/user/current      --> github.com/fabric-app/controller/api/v1.CurrentUser (5 handlers)
[GIN-debug] GET    /api/v1/user/refresh      --> github.com/fabric-app/controller/api/v1.RefreshToken (5 handlers)
[GIN-debug] POST   /api/v1/user/logout       --> github.com/fabric-app/controller/api/v1.Logout (5 handlers)
[GIN-debug] POST   /api/v1/user/password     --> github.com/fabric-app/controller/api/v1.Password (5 handlers)
[GIN-debug] POST   /api/v1/user/update       --> github.com/fabric-app/controller/api/v1.ModifyUser (5 handlers)
[GIN-debug] POST   /api/v1/user/record       --> github.com/fabric-app/controller/api/v1.Record (5 handlers)
[GIN-debug] GET    /api/v1/user/operType     --> github.com/fabric-app/controller/api/v1.Operations (5 handlers)
[GIN-debug] POST   /api/v1/user/setHeader    --> github.com/fabric-app/controller/api/v1.SetHeader (5 handlers)
[GIN-debug] POST   /api/v1/user/getHeader    --> github.com/fabric-app/controller/api/v1.GetHeader (5 handlers)
[GIN-debug] POST   /api/v1/user/revoke       --> github.com/fabric-app/controller/api/v1.Revoker (5 handlers)
[GIN-debug] GET    /api/v1/bcs/info          --> github.com/fabric-app/controller/api/v1.BcInfo (5 handlers)
[GIN-debug] POST   /api/v1/bcs/transactions  --> github.com/fabric-app/controller/api/v1.Transactions (5 handlers)
[GIN-debug] GET    /api/v1/bcs/points        --> github.com/fabric-app/controller/api/v1.Points (5 handlers)
[GIN-debug] POST   /api/v1/trace/sensor      --> github.com/fabric-app/controller/api/v1.Sensors (5 handlers)
[GIN-debug] POST   /api/v1/trace/picture     --> github.com/fabric-app/controller/api/v1.Pictures (5 handlers)
[GIN-debug] POST   /api/v1/trace/farmData    --> github.com/fabric-app/controller/api/v1.Farms (5 handlers)
[GIN-debug] POST   /api/v1/trace/verify      --> github.com/fabric-app/controller/api/v1.Verifier (5 handlers)
[GIN-debug] POST   /api/v1/trace/upload      --> github.com/fabric-app/controller/api/v1.Uploader (5 handlers)


```
### Swaggo

> http://127.0.0.1:8000/swagger/index.html

## 项目结构概览
```
├── config 配置文件,包含fabric连接配置、web server参数
│   └── crypto-config fabric证书材料
│   ├── api.ini Web server 配置 
├── docs：api文档swagger    
│   └── sql：sql执行语句  
├── middleware：中间件
│   └── jwt：认证中间件
├── model：引用数据库模型、fabric sdk存储层
│   └── bcs: fabric blockchain服务
├── pkg：第三方包和公共模块
│   ├── app：gin engine
│   ├── e： 错误编码和错误信息
│   ├── logging：日志模块
│   ├── setting：go-ini包
│   └── util：工具库 
└── routers：路由处理
│    └── api：controller 逻辑梳理
│        └── v1：controller逻辑处理 
└── service：逻辑处理
└── runtime：日志，文件缓存存放
└── test：单元测试
│   ├── chaincode 链码
│   ├── simple-network fabric 本地测试网络
│   └── header 用户头像
└── main.go：入口文件 
```






