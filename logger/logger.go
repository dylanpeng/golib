package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

type Config struct {
	FilePath   string `toml:"file_path" json:"file_path"`
	Level      string `toml:"level" json:"level"`
	TimeFormat string `toml:"time_format" json:"time_format"`
	MaxAgeDay  int    `toml:"max_age_day" json:"max_age_day"`
}

func DefaultConfig() *Config {
	return &Config{
		FilePath:   "./logs/project",
		Level:      "debug",
		TimeFormat: "2006-01-02 15:04:05.000",
		MaxAgeDay:  30,
	}
}

type Logger struct {
	conf   *Config
	sugar  *zap.SugaredLogger
	writer io.WriteCloser
}

func (l *Logger) Init() error {
	if l.conf == nil {
		l.conf = DefaultConfig()
	}

	// 日志编码默认配置
	encoderConf := zap.NewDevelopmentEncoderConfig()
	// 日志时间格式
	encoderConf.EncodeTime = zapcore.TimeEncoderOfLayout(l.conf.TimeFormat)
	// 按天分文件
	loggerWriter := getWriter(l.conf.FilePath, l.conf.MaxAgeDay)
	// set logger writer
	l.writer = loggerWriter
	// convert zap writer
	writer := zapcore.AddSync(loggerWriter)
	level, err := zapcore.ParseLevel(l.conf.Level)

	if err != nil {
		return err
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConf),
		//zapcore.NewMultiWriteSyncer(writer, zapcore.AddSync(os.Stdout)),
		writer,
		level,
	)

	// warn error级别日志打出调用链,并跳过封装路径
	p := zap.New(core, zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller(), zap.AddCallerSkip(1))
	l.sugar = p.Sugar()

	return nil
}

func (l *Logger) GetWriter() io.Writer {
	return l.writer
}

func (l *Logger) Debug(args ...any) {
	if l == nil {
		return
	}

	l.sugar.Debugln(args...)
}

func (l *Logger) Debugf(format string, args ...any) {
	if l == nil {
		return
	}

	l.sugar.Debugf(format, args...)
}

func (l *Logger) Info(args ...any) {
	if l == nil {
		return
	}

	l.sugar.Infoln(args...)
}

func (l *Logger) Infof(format string, args ...any) {
	if l == nil {
		return
	}

	l.sugar.Infof(format, args...)
}

func (l *Logger) Warn(args ...any) {
	if l == nil {
		return
	}

	l.sugar.Warnln(args...)
}

func (l *Logger) Warningf(format string, args ...any) {
	if l == nil {
		return
	}

	l.sugar.Warnf(format, args...)
}

func (l *Logger) Error(args ...any) {
	if l == nil {
		return
	}

	l.sugar.Errorln(args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	if l == nil {
		return
	}

	l.sugar.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...any) {
	if l == nil {
		return
	}

	l.sugar.Fatalln(args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	if l == nil {
		return
	}

	l.sugar.Fatalf(format, args...)
}

func (l *Logger) Sync() error {
	return l.sugar.Sync()
}

func (l *Logger) Close() error {
	return l.writer.Close()
}

func NewLogger(conf *Config) (l *Logger, err error) {
	l = &Logger{conf: conf}
	err = l.Init()

	return
}

func getWriter(filename string, maxAgeDay int) io.WriteCloser {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		// 没有使用go风格反人类的format格式
		filename+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxAgeDay)),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
