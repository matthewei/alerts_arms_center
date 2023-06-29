package zaplogger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strconv"
)

const loggerKey = iota

var Logger *zap.Logger

type Configurations struct {
	LogLevel          string `yaml:"log_level"`            // 日志打印级别 debug  info  warning  error
	LogFormat         string `yaml:"log_format"`           // 输出日志格式	logfmt, json
	LogPath           string `yaml:"log_path"`             // 输出日志文件路径
	LogFileName       string `yaml:"log_file_name"`        // 输出日志文件名称
	LogFileMaxSize    int    `yaml:"log_file_max_size"`    // 【日志分割】单个日志文件最多存储量 单位(mb)
	LogFileMaxBackups int    `yaml:"log_file_max_backups"` // 【日志分割】日志备份文件最多数量
	LogMaxAge         int    `yaml:"log_max_age"`          // 日志保留时间，单位: 天 (day)
	LogCompress       bool   `yaml:"log_compress"`         // 是否压缩日志
	LogStdout         bool   `yaml:"log_stdout"`           // 是否输出到控制台
}

type ZapLogger struct {
	*zap.SugaredLogger
}

func New(configuration Configurations) (*ZapLogger, error) {
	logLevel := map[string]zapcore.Level{
		"debug":   zapcore.DebugLevel,
		"info":    zapcore.InfoLevel,
		"warning": zapcore.WarnLevel,
		"error":   zapcore.ErrorLevel,
	}
	writeSyncer, err := getLogWriter(configuration) // 日志文件配置 文件位置和切割
	if err != nil {
		return nil, err
	}
	encoder := getEncoder(configuration)          // 获取日志输出编码
	level, ok := logLevel[configuration.LogLevel] // 日志打印级别
	if !ok {
		level = logLevel["info"]
	}
	logCore := zapcore.NewCore(encoder, writeSyncer, level)

	//创建具体的Logger
	core := zapcore.NewTee(logCore)
	logger := zap.New(core, zap.AddCaller())
	zapLogger := new(ZapLogger)
	zapLogger.SugaredLogger = logger.Sugar()
	defer zapLogger.Sync()
	Logger = logger
	return zapLogger, nil
}

func getLogWriter(config Configurations) (zapcore.WriteSyncer, error) {
	// 判断日志路径是否存在，如果不存在就创建
	if exist := IsExist(config.LogPath); !exist {
		if config.LogPath == "" {
			config.LogPath = getCurrentDirectory()
		}
		if err := os.MkdirAll(config.LogPath, os.ModePerm); err != nil {
			config.LogPath = getCurrentDirectory()
			if err := os.MkdirAll(config.LogPath, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}
	if config.LogFileName == "" {
		config.LogFileName = "alerts_arms_center.log"
	}
	// 日志文件 与 日志切割 配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(config.LogPath, config.LogFileName), // 日志文件路径
		MaxSize:    config.LogFileMaxSize,                             // 单个日志文件最大多少 mb
		MaxBackups: config.LogFileMaxBackups,                          // 日志备份数量
		MaxAge:     config.LogMaxAge,                                  // 日志最长保留时间
		Compress:   config.LogCompress,                                // 是否压缩日志
	}
	if config.LogStdout {
		// 日志同时输出到控制台和日志文件中
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout)), nil
	} else {
		// 日志只输出到日志文件
		return zapcore.AddSync(lumberJackLogger), nil
	}
}

// IsExist 判断文件或者目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 获取当前目录
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(filepath.Join(".")))
	if err != nil {
		panic(err)
	}
	logFileDir := filepath.Join(dir, "/logs/")
	return logFileDir
}

// Encoder
func getEncoder(config Configurations) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // log 时间格式 例如: 2021-09-11t20:05:54.852+0800
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 输出level序列化为全大写字符串，如 INFO DEBUG ERROR
	if config.LogFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig) // 以json格式写入
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 以logfmt格式写入
}

// 给指定的context添加字段（关键方法）
func NewContext(ctx *gin.Context, fields ...zapcore.Field) {
	ctx.Set(strconv.Itoa(loggerKey), WithContext(ctx).With(fields...))
}

// 从指定的context返回一个zap实例（关键方法）
func WithContext(ctx *gin.Context) *zap.Logger {
	if ctx == nil {
		return Logger
	}
	l, _ := ctx.Get(strconv.Itoa(loggerKey))
	ctxLogger, ok := l.(*zap.Logger)
	if ok {
		return ctxLogger
	}
	return Logger
}
