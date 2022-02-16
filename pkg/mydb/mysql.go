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
func CreateTables() {
	createArticlesSQL := `SET NAMES utf8mb4;`
	MysqlClient.Exec(createArticlesSQL)
	createArticlesSQL = `SET time_zone = '+00:00';`
	MysqlClient.Exec(createArticlesSQL)
	createArticlesSQL = `SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';`
	MysqlClient.Exec(createArticlesSQL)
	createArticlesSQL = `DROP DATABASE bytecamp;`
	MysqlClient.Exec(createArticlesSQL)
	createArticlesSQL = `CREATE DATABASE bytecamp;`
	MysqlClient.Exec(createArticlesSQL)

	createArticlesSQL = `USE bytecamp;`
	MysqlClient.Exec(createArticlesSQL)

	createArticlesSQL = `CREATE TABLE users(
  user_id bigint NOT NULL AUTO_INCREMENT,
  nickname varchar(20) NOT NULL,
  username varchar(20) NOT NULL,
  password varchar(100) NOT NULL,
  user_type tinyint NOT NULL,
  deleted_at datetime,
  PRIMARY KEY (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	MysqlClient.Exec(createArticlesSQL)

	createArticlesSQL = `CREATE TABLE courses (
  course_id bigint NOT NULL AUTO_INCREMENT,
  name varchar(20) NOT NULL,
  cap int NOT NULL,
  teacher_id bigint,
  deleted_at datetime,
  PRIMARY KEY (course_id),
  FOREIGN KEY (teacher_id) REFERENCES users(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	MysqlClient.Exec(createArticlesSQL)

	createArticlesSQL = `CREATE TABLE bind_courses (
  teacher_id bigint NOT NULL,
  course_id bigint NOT NULL,
  deleted_at datetime,
  FOREIGN KEY (teacher_id) REFERENCES users(user_id),
  FOREIGN KEY (course_id) REFERENCES courses(course_id),
  PRIMARY KEY (teacher_id, course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	MysqlClient.Exec(createArticlesSQL)

	createArticlesSQL = `CREATE TABLE select_courses (
  student_id bigint NOT NULL,
  course_id bigint NOT NULL,
  deleted_at datetime,
  FOREIGN KEY (student_id) REFERENCES users(user_id),
  FOREIGN KEY (course_id) REFERENCES courses(course_id),
  PRIMARY KEY (student_id, course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	MysqlClient.Exec(createArticlesSQL)

	createArticlesSQL = `INSERT INTO users (nickname, username, password, user_type) 
VALUES ('JudgeAdmin', 'JudgeAdmin', 'JudgePassword2022', 1);`
	MysqlClient.Exec(createArticlesSQL)

}
