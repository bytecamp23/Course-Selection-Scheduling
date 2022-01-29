## 分工

zwm、dzh、lyh

- 登录

- 成员

- 排课

yzg

- 文档
- sql、redis包装
- 排课求解器

zlf

- 抢课

## 数据库设计

| users      |             |                                       | 用户表             |
| ---------- | ----------- | ------------------------------------- | ------------------ |
| user_id    | bigint      | PRIMARY KEY, AUTO_INCREMENT, NOT NULL |                    |
| nickname   | varchar(20) | NOT NULL                              |                    |
| username   | varchar(20) | NOT NULL                              |                    |
| password   | varchar(20) | NOT NULL                              |                    |
| user_type  | tinyint     | NOT NULL                              |                    |
| deleted_at | datetime    |                                       | gorm自动支持软删除 |

> Q：请问老师和课程的关系是多对多吗？就是A、B、C、D老师都可能上高数课，然后A除了高数课还有线性代数、概率论等课。文档里写的是一个老师对应一个课程，是只有A上高数课，且高数课只有A老师可以选择吗？
>
> A：现实场景是有的，不过我们这里做了限制。一个课程只能绑定一个老师，一个老师可以有多个课程。
>
> 对于算法题，老师和课程是一对一关系。             

| courses    |             |                                       | 课表                 |
| ---------- | ----------- | ------------------------------------- | -------------------- |
| course_id  | bigint      | PRIMARY KEY, AUTO_INCREMENT, NOT NULL |                      |
| name       | varchar(20) | NOT NULL                              |                      |
| cap        | int         | NOT NULL                              |                      |
| teacher_id | bigint      | FOREIGN KEY REFERENCES users(user_id) | 排课后 老师:课程=1:1 |
| deleted_at | datetime    |                                       | gorm自动支持软删除   |



| bind_courses |          |                                                              | 教师绑定课表         |
| ------------ | -------- | ------------------------------------------------------------ | -------------------- |
| teacher_id   | bigint   | PRIMARY KEY, FOREIGN KEY REFERENCES users(user_id), NOT NULL |                      |
| course_id    | bigint   | PRIMARY KEY, FOREIGN KEY REFERENCES courses(course_id), NOT NULL | 排课前 老师:课程=1:n |
| deleted_at   | datetime |                                                              | gorm自动支持软删除   |



| select_courses |          |                                                              | 学生选课表         |
| -------------- | -------- | ------------------------------------------------------------ | ------------------ |
| student_id     | bigint   | PRIMARY KEY, FOREIGN KEY REFERENCES users(user_id), NOT NULL |                    |
| course_id      | bigint   | PRIMARY KEY, FOREIGN KEY REFERENCES courses(course_id), NOT NULL | 学生:课程=n:m      |
| deleted_at     | datetime |                                                              | gorm自动支持软删除 |



```sql
SET NAMES utf8mb4;
SET time_zone = '+00:00';
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP DATABASE `bytecamp`;
CREATE DATABASE `bytecamp`;

USE `bytecamp`;

CREATE TABLE `users` (
  `user_id` bigint NOT NULL AUTO_INCREMENT,
  `nickname` varchar(20) NOT NULL,
  `username` varchar(20) NOT NULL,
  `password` varchar(100) NOT NULL,
  `user_type` tinyint NOT NULL,
  `deleted_at` datetime,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `courses` (
  `course_id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL,
  `cap` int NOT NULL,
  `teacher_id` bigint,
  `deleted_at` datetime,
  PRIMARY KEY (`course_id`),
  FOREIGN KEY (`teacher_id`) REFERENCES users(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `bind_courses` (
  `teacher_id` bigint NOT NULL,
  `course_id` bigint NOT NULL,
  `deleted_at` datetime,
  FOREIGN KEY (`teacher_id`) REFERENCES users(user_id),
  FOREIGN KEY (`course_id`) REFERENCES courses(course_id),
  PRIMARY KEY (`teacher_id`, `course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `select_courses` (
  `student_id` bigint NOT NULL,
  `course_id` bigint NOT NULL,
  `deleted_at` datetime,
  FOREIGN KEY (`student_id`) REFERENCES users(user_id),
  FOREIGN KEY (`course_id`) REFERENCES courses(course_id),
  PRIMARY KEY (`student_id`, `course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `users` (`nickname`, `username`, `password`, `user_type`) 
VALUES ('admin', 'admin123', '123456', 1); #默认管理员账号 做示例 实际不存明文
```

