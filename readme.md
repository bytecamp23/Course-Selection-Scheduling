# Course-Selection-Scheduling

## 概要

该系统为基于Go语言的选排课系统，具有登录、成员、排课、抢课模块。

### 模块及考察点

- 登录模块
  - 考察登录的设计与实现，对 HTTP 协议的理解。
    - 账密登录
    - Cookie Session

- 成员模块
  - 考察工程实现能力。
    - CURD 及对数据库的操作
    - 参数校验
      - 参数长度
      - 弱密码校验
    - 权限判断

- 排课模块
  - 主要考察算法（二分图匹配）的实现。

- 抢课模块
  - 主要考察简单秒杀场景的设计。

具体需求：[课程1 for 营员 - 飞书云文档 (feishu.cn)](https://bytedance.feishu.cn/docx/doxcnVuJvnjMgf99tiGE3fSiIBP#doxcnWUQ4syyOOG6qCKL2XnUbmP)

### 分工

- 张伟明、邓泽晖、刘曜辉负责登录、成员、排课模块（排课求解器除外）。
- 杨志刚负责文档编写、基础架构、排课求解器。
- 曾林藩负责抢课模块。



## 系统设计与实现

### 系统技术框架

- 主体语言: Golang 

- WEB框架: Gin

- 持久层框架: Gorm、RediGo 

- 关系型数据库: Mysql 

- 缓存数据库: Redis

### 目录结构

```
.
└── config //配置目录
│   └── dev //具体环境
│       ├── log.yml //日志配置
│       ├── mysql.yml //mysql配置
│       ├── redis.yml //redis配置
│       ├── server.yml //httpserver配置
│       └── session.yml //session配置
├── go.mod
├── go.sum
├── internal //内部目录
│   ├── api //api目录 根据大作业介绍-接口设计 controller + model
│   │   ├── auth
│   │   │   ├── controller.go
│   │   │   └── model.go
│   │   ├── course
│   │   │   ├── controller.go
│   │   │   └── model.go
│   │   ├── member
│   │   │   ├── controller.go
│   │   │   └── model.go
│   │   ├── student
│   │   │   ├── controller.go
│   │   │   └── model.go
│   │   └── teacher
│   │       ├── controller.go
│   │       └── model.go
│   ├── global// 全局目录
│   │   ├── global.go
│   │   └── types.go 
│   └── server// httpserver目录
│       ├── router.go
│       └── server.go
├── logs //日志目录
│   └── dev //具体环境
├── main.go 
├── pkg //外部包目录
│   ├── config //配置文件定义目录
│   │   └── config.go
│   ├── mydb //mysql定义目录
│   │   └── mysql.go
│   └── myredis //redis定义目录
│       └── redis.go
├── readme.md
└── utils //工具目录
    └── loadCfg.go //加载配置文件
```

### 数据库设计

- users 用户信息表

|     键     |  数据类型   |                 约束                  |   注释   |
| :--------: | :---------: | :-----------------------------------: | :------: |
|  user_id   |   bigint    | PRIMARY KEY, AUTO_INCREMENT, NOT NULL |  用户ID  |
|  nickname  | varchar(20) |               NOT NULL                |   昵称   |
|  username  | varchar(20) |               NOT NULL                |  用户名  |
|  password  | varchar(20) |               NOT NULL                |   密码   |
| user_type  |   tinyint   |               NOT NULL                | 用户类型 |
| deleted_at |  datetime   |                                       |  软删除  |

- courses 课程表 排课后老师:课程=1:1

|     键     |  数据类型   |                 约束                  |   注释   |
| :--------: | :---------: | :-----------------------------------: | :------: |
| course_id  |   bigint    | PRIMARY KEY, AUTO_INCREMENT, NOT NULL |  课程ID  |
|    name    | varchar(20) |               NOT NULL                |  课程名  |
|    cap     |     int     |               NOT NULL                | 课程容量 |
| teacher_id |   bigint    | FOREIGN KEY REFERENCES users(user_id) |  教师ID  |
| deleted_at |  datetime   |                                       |  软删除  |

- bind_courses 教师绑定课程表 排课前 老师:课程=1:n

|     键     | 数据类型 |                             约束                             |  注释  |
| :--------: | :------: | :----------------------------------------------------------: | :----: |
| teacher_id |  bigint  | PRIMARY KEY, FOREIGN KEY REFERENCES users(user_id), NOT NULL | 教师ID |
| course_id  |  bigint  | PRIMARY KEY, FOREIGN KEY REFERENCES courses(course_id), NOT NULL | 课程ID |
| deleted_at | datetime |                                                              | 软删除 |

- select_courses 学生选课表 学生:课程=n:m

|     键     | 数据类型 |                             约束                             |  注释  |
| :--------: | :------: | :----------------------------------------------------------: | :----: |
| student_id |  bigint  | PRIMARY KEY, FOREIGN KEY REFERENCES users(user_id), NOT NULL | 学生ID |
| course_id  |  bigint  | PRIMARY KEY, FOREIGN KEY REFERENCES courses(course_id), NOT NULL | 课程ID |
| deleted_at | datetime |                                                              | 软删除 |

- 数据库脚本

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
VALUES ('JudgeAdmin', 'JudgeAdmin', 'JudgePassword2022', 1); #默认管理员账号 做示例 实际不存明文
```



