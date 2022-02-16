package mydb

import (
	"Course-Selection-Scheduling/types"
	"Course-Selection-Scheduling/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var MysqlClient *gorm.DB

type User struct {
	UserId    string `gorm:"primaryKey;autoIncrement"`
	Nickname  string
	Username  string
	Password  string
	UserType  types.UserType
	DeletedAt gorm.DeletedAt
}

type Course struct {
	CourseId  string `gorm:"primaryKey;autoIncrement"`
	Name      string
	Cap       int
	TeacherId *string
	DeletedAt gorm.DeletedAt
}

type BindCourse struct {
	TeacherId string `gorm:"primaryKey"`
	CourseId  string `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt
}

type SelectCourse struct {
	StudentId string `gorm:"primaryKey"`
	CourseId  string `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt
}

func NewMysqlConn(cfg *utils.Mysql) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("mysql open error! " + err.Error())
	}
	db, _ := mysqlDb.DB()
	// 设置最大连接数
	db.SetMaxOpenConns(100)
	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)
	return mysqlDb
}
