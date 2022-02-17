package mydb

import (
	"Course-Selection-Scheduling/pkg/myredis"
)

func ClearDB() {
	db := MysqlClient
	createArticlesSQL := `SET NAMES utf8mb4;`
	db.Exec(createArticlesSQL)
	createArticlesSQL = `SET time_zone = '+00:00';`
	db.Exec(createArticlesSQL)
	createArticlesSQL = `SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';`
	db.Exec(createArticlesSQL)
	createArticlesSQL = `DROP DATABASE bytecamp;`
	db.Exec(createArticlesSQL)
	createArticlesSQL = `CREATE DATABASE bytecamp;`
	db.Exec(createArticlesSQL)

	createArticlesSQL = `USE bytecamp;`
	db.Exec(createArticlesSQL)

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
	db.Exec(createArticlesSQL)

	createArticlesSQL = `CREATE TABLE courses (
  course_id bigint NOT NULL AUTO_INCREMENT,
  name varchar(20) NOT NULL,
  cap int NOT NULL,
  teacher_id bigint,
  deleted_at datetime,
  PRIMARY KEY (course_id),
  FOREIGN KEY (teacher_id) REFERENCES users(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	db.Exec(createArticlesSQL)

	createArticlesSQL = `CREATE TABLE bind_courses (
  teacher_id bigint NOT NULL,
  course_id bigint NOT NULL,
  deleted_at datetime,
  FOREIGN KEY (teacher_id) REFERENCES users(user_id),
  FOREIGN KEY (course_id) REFERENCES courses(course_id),
  PRIMARY KEY (teacher_id, course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	db.Exec(createArticlesSQL)

	createArticlesSQL = `CREATE TABLE select_courses (
  student_id bigint NOT NULL,
  course_id bigint NOT NULL,
  deleted_at datetime,
  FOREIGN KEY (student_id) REFERENCES users(user_id),
  FOREIGN KEY (course_id) REFERENCES courses(course_id),
  PRIMARY KEY (student_id, course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	db.Exec(createArticlesSQL)

	createArticlesSQL = `INSERT INTO users (nickname, username, password, user_type) 
VALUES ('JudgeAdmin', 'JudgeAdmin', 'JudgePassword2022', 1);`
	db.Exec(createArticlesSQL)

	myredis.Flushdb()
}
