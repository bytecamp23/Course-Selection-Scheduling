# Course-Selection-Scheduling

[TOC]

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

- 张伟明负责登录和成员的接口、压力测试。
- 邓泽晖负责排课模块（排课求解器除外）、请求参数校验器。
- 杨志刚负责文档编写、基础架构、排课求解器、逻辑测试。
- 曾林藩负责抢课模块。



## 系统设计与实现

### 系统技术框架

- 主体语言: Golang 

- WEB框架: Gin

- 持久层框架: Gorm、RediGo 

- 关系型数据库: Mysql 

- 缓存数据库: Redis

### 目录结构

参照MVC架构

```
.
├── api//存放业务代码
│   ├── controllers//控制器
│   │   ├── auth
│   │   │   └── auth.go
│   │   ├── course
│   │   │   └── course.go
│   │   ├── member
│   │   │   └── member.go
│   │   ├── student
│   │   │   └── student.go
│   │   └── teacher
│   │       └── teacher.go
│   └── models//数据模型
│       ├── auth
│       │   └── auth.go
│       ├── course
│       │   ├── course.go
│       │   └── matchSolve.go
│       ├── member
│       │   └── member.go
│       ├── student
│       │   └── student.go
│       └── teacher
│           └── teacher.go
├── config//存放配置文件
│   └── dev
│       ├── log.yml
│       ├── mysql.yml
│       ├── rabbitMQ.yml
│       ├── redis.yml
│       ├── server.yml
│       └── session.yml
├── go.mod
├── go.sum
├── logs//存放日志文件
│   └── dev
├── main.go
├── pkg//存放底层代码
│   ├── logger
│   │   └── logger.go
│   ├── mydb
│   │   └── mysql.go
│   ├── myredis
│   │   └── redis.go
│   ├── rabbitmq
│   │   └── rabbitmq.go
│   └── server
│       ├── router.go
│       └── server.go
├── readme.md
├── types//一些定义
│   └── types.go
└── utils//存放工具
    ├── loadConfig.go
    ├── md5.go
    ├── min.go
    └── validator.go
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
VALUES ('JudgeAdmin', 'JudgeAdmin', 'JudgePassword2022', 1); 
```

### 排课求解器实现

本质为求解二分图最大匹配并输出方案，通常采用匈牙利（Hungarian）或Dinic算法解决。

在点数为 $N$、边数为 $M$ 的二分图上：

- 临接表实现的匈牙利算法的最劣时间复杂度为$O(MN)$，空间复杂度为$O(N+M)$
- 临接表实现的Dinic算法的最劣时间复杂度为 $O(M\sqrt N)$，空间复杂度为$O(N+M)$

由于给定的ID范围不一定从$0$连续递增，我们先采取离散化算法，将ID离散到$[0,N]$的值域上，降低空间复杂度。

由于Dinic算法在时间和空间上的常数较大，小数据量下往往不如匈牙利算法优秀，本项目采取根据数据量分别采用两种方法。

具体而言当$N*M>1000000$我们采用Dinic算法，否则我们采用匈牙利算法。

### 异常处理

#### 用户模块

* create

ParamInvalid
UserHasExisted

* get

ParamInvalid
UserHasDeleted
UserNotExisted

* list

ParamInvalid

* update

ParamInvalid
UserHasDeleted
UserNotExisted

* delete

ParamInvalid
UserHasDeleted
UserNotExisted

#### 登陆模块

* login

ParamInvalid
UserHasDeleted
UserNotExisted
WrongPassword

* logout

ParamInvalid
LoginRequired

* whoami

ParamInvalid
LoginRequired

#### 排课模块

* create

ParamInvalid
课程已存在？UnknownError

* get

ParamInvalid
CourseNotExisted

* schedule

ParamInvalid

UnknownError

* bind

ParamInvalid
CourseHasBound
PermDenied
CourseNotExisted

* unbind

ParamInvalid
CourseNotBound
PermDenied
CourseNotExisted

* get

ParamInvalid
CourseNotBound
CourseNotExisted

#### 选课模块

* book_course

ParamInvalid

CourseNotAvailable

StudentNotExisted

CourseNotExisted

StudentHasCourse

RepeatRequest

UnknownError

* get

ParamInvalid

StudentNotExisted

StudentHasNoCourse

### 抢课模块实现

#### 抢课

核心流程：

![2A1622099B733158E5B1885336852F83](/Users/lmessi/OneDrive/Project/Course-Selection-Scheduling/assets/2A1622099B733158E5B1885336852F83.jpg)

判断非法课程号和学生号：redis辅助实现

查课表：redis辅助实现

