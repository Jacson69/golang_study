package main

// 日志切割
//func newLfsHook(logLevel *string, maxRemainCnt uint) log.Hook {
//	writer, err := rotatelogs.New(
//		"logName"+".%Y%m%d%H",
//		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
//		rotatelogs.WithLinkName("logName"),
//
//		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
//		rotatelogs.WithRotationTime(time.Hour),
//
//		// WithMaxAge和WithRotationCount二者只能设置一个，
//		// WithMaxAge设置文件清理前的最长保存时间，
//		// WithRotationCount设置文件清理前最多保存的个数。
//		//rotatelogs.WithMaxAge(time.Hour*24),
//		rotatelogs.WithRotationCount(maxRemainCnt),
//	)
//
//	if err != nil {
//		log.Errorf("config local file system for logger error: %v", err)
//	}
//
//	level, ok := logLevels[*logLevel]
//
//	if ok {
//		log.SetLevel(level)
//	} else {
//		log.SetLevel(log.WarnLevel)
//	}
//
//	lfsHook := lfshook.NewHook(lfshook.WriterMap{
//		log.DebugLevel: writer,
//		log.InfoLevel:  writer,
//		log.WarnLevel:  writer,
//		log.ErrorLevel: writer,
//		log.FatalLevel: writer,
//		log.PanicLevel: writer,
//	}, &log.TextFormatter{DisableColors: true})
//
//	return lfsHook
//}
import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"time"
)

var (
	logFilePath = "./" //文件存储路径
	logFileName = "system.log"
)

func main() {
	r := gin.Default()
	//添加中间件，主要实现log日志的生成
	r.Use(logMiddleware())
	r.GET("/logrus2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "响应成功",
			"data": "ok",
		})
	})
	r.Run(":9090")
}
func logMiddleware() gin.HandlerFunc {
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //创建一个log日志文件
	if err != nil {
		fmt.Println(err)
	}
	//实例化
	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	logger.Out = file //设置log的默认文件输出
	logWriter, err := rotatelogs.New(
		//分割后的文件名称
		fileName+".%Y%m%d.log",
		//生成软链接，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		//设置最大保存时间7天
		rotatelogs.WithMaxAge(7*24*time.Hour),
		//设置日志切割时间间隔（1天）
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	//hook机制的设置
	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	//给logrus添加hook
	logger.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
	return func(c *gin.Context) {
		c.Next()
		//请求方法
		method := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIp := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"client_ip":   clientIp,
			"req_method":  method,
			"req_url":     reqUrl,
		}).Info()
	}
}
