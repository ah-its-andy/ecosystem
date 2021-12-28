package logging

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logrus "github.com/sirupsen/logrus"
)

type Logger interface {
	Error(err error)
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type smartLogger struct {
	serviceName string
	dirPath     string
	fileName    string
	logWriter   *rotatelogs.RotateLogs
}

func NewSmartLogger(serviceName string, dirPath string) Logger {
	r := &smartLogger{
		serviceName: serviceName,
		dirPath:     dirPath,
	}
	r.fileName = fmt.Sprintf("%s.log", serviceName)
	path := filepath.Join(r.dirPath, r.fileName)
	logWriter, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		//保留最近30天的日志
		rotatelogs.WithMaxAge(time.Duration(720)*time.Hour),
		//每30分钟分割一个文件
		rotatelogs.WithRotationTime(time.Duration(30)*time.Minute),
	)
	if err != nil {
		panic(errors.Unwrap(err))
	}
	r.logWriter = logWriter
	return r
}

func (logger *smartLogger) Error(err error) {
	logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName}).Error(err)
}

func (logger *smartLogger) Errorf(format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName}).Errorf(format, args...)
}

func (logger *smartLogger) Debugf(format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName}).Debugf(format, args...)
}
