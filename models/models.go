package models

import (
	"fmt"
	"github.com/8treenet/gcache"
	"github.com/fabric-app/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var db *gorm.DB

//基础模型 都用Id 做自增键  都有创建时间和更新时间
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

//初始化数据库连接池设置
func init() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)
	//从设置中取设置
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	//赋值数据连接变量
	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		log.Println(err)
	}

	opt := gcache.DefaultOption{}
	opt.Expires = 300              //缓存时间，默认60秒。范围 30-900
	opt.Level = gcache.LevelSearch //缓存级别，默认LevelSearch。LevelDisable:关闭缓存，LevelModel:模型缓存， LevelSearch:查询缓存
	opt.AsyncWrite = true         //异步缓存更新, 默认false。 insert update delete 成功后是否异步更新缓存
	opt.PenetrationSafe = false    //开启防穿透, 默认false。

	//缓存中间件 注入到Gorm
	gcache.AttachDB(db, &opt, &gcache.RedisOption{Addr: "localhost:6379"})

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	db.SingularTable(false)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	//createTables(db)

}

func createTables(db *gorm.DB) {
	//if db.HasTable(&User{}) {  // 创建 User 表
	//	db.AutoMigrate(&User{})
	//}else {
	//	db.CreateTable(&User{})
	//}
	db.AutoMigrate(&User{}, &Transaction{}, &FarmType{})
}

func CloseDB() {
	defer db.Close()
}
