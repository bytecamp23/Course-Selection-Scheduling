package global

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type User struct {
	UserId    int64 `gorm:"primaryKey"`
	Nickname  string
	Username  string
	Password  string
	UserType  UserType
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

var MysqlClient *gorm.DB

func NewMysqlConn(name, password, database string) {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", name, password, database)
	mysqlClient, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("mysql open error! " + err.Error())
	}
	MysqlClient = mysqlClient
}
