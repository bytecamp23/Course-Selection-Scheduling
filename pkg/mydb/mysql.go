package mydb

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	UserId    int64 `gorm:"primaryKey"`
	Nickname  string
	Username  string
	Password  string
	UserType  global.UserType
	DeletedAt gorm.DeletedAt
}

type Course struct {
	CourseId  int64 `gorm:"primaryKey"`
	Name      string
	Cap       int
	TeacherId *int64
	DeletedAt gorm.DeletedAt
}

type BindCourse struct {
	TeacherId int64 `gorm:"primaryKey"`
	CourseId  int64 `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt
}

type SelectCourse struct {
	StudentId int64 `gorm:"primaryKey"`
	CourseId  int64 `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt
}

func NewMysqlConn(cfg *config.Mysql) *gorm.DB {
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
	return mysqlDb
}
