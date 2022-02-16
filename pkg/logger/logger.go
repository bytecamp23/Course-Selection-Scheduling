package logger

import (
	"fmt"
	"log"
)

func FilePanic(path string, err error) {
	if err != nil {
		log.Panic(fmt.Sprintf("read config file failed, path: %v. err: %v", path, err.Error()))
	}
}
func ParsePanic(path string, err error) {
	if err != nil {
		log.Panic(fmt.Sprintf("parse config file failed, path: %v. err: %v", path, err.Error()))
	}
}
