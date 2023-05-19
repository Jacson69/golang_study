package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// 设置日志格式为json格式
	log.SetFormatter(&log.JSONFormatter{})

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	log.SetLevel(log.WarnLevel)
}

func main() {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")

	requestLogger := log.WithFields(log.Fields{"request_id": "2145643523412", "user_ip": "123"})
	requestLogger.Info("something happened on that request")
	requestLogger.Warn("something not great happened")

}

// logrus提供了New()函数来创建一个logrus的实例。
// 项目中，可以创建任意数量的logrus实例。
//var log = logrus.New()
//
//func main() {
//	// 为当前logrus实例设置消息的输出，同样地，
//	// 可以设置logrus实例的输出到任意io.writer
//	log.Out = os.Stdout
//
//	// 为当前logrus实例设置消息输出格式为json格式。
//	// 同样地，也可以单独为某个logrus实例设置日志级别和hook，这里不详细叙述。
//	log.Formatter = &logrus.JSONFormatter{}
//
//	log.WithFields(logrus.Fields{
//		"animal": "walrus",
//		"size":   10,
//	}).Info("A group of walrus emerges from the ocean")
//}
