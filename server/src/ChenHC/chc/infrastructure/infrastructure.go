package infrastructure

import (
	"ChenHC/chc/cache"
	"ChenHC/chc/db"
	"ChenHC/chc/logger_core"
	"ChenHC/chc/option"
	"ChenHC/chc/session"
	"ChenHC/chc/wechat"
	"ChenHC/utils"
	"database/sql"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/asdine/storm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"fmt"
)
//基础设施
type Infrastructure struct {
	opts           atomic.Value //并发并行的原子操作
	Logger         *zap.Logger //logger包
	DB             *sql.DB   //数据库
	stormDB        *storm.DB  //stormDb
	SessionManager *session.SessionManager //seesion管理器
	CacheManager   *cache.CacheManager  //cache stormdb
}

func NewInfrastructure(opts *option.Options) (i *Infrastructure) {
	proDir, err := utils.GetProDir()
	if err != nil {
		log.Panicln("get prodir err: " + err.Error())
	}

	i = &Infrastructure{}
	i.setOpts(opts)
	i.DB = db.NewDB(opts.DBHost, opts.DBPort, opts.DBUser, opts.DBPassword, opts.DBName, opts.DBDriver)
	i.CacheManager = cache.NewCacheManager(proDir + "/my.db")
	i.SessionManager = new(session.SessionManager)
	i.SessionManager.Init(
		opts.SessionKey,
		opts.MaxLifeTime,
		opts.Path,
		opts.HTTPOnly,
		opts.MaxAge,
		i.CacheManager,
		opts.BuildType)
	i.InitLogger()

	wechat.InitWeChat( opts.AppID, opts.AppSecret, i.CacheManager)
	fmt.Println("初始化Infrastructure完成。")
	return
}

func (i *Infrastructure) GetOpts() *option.Options {
	return i.opts.Load().(*option.Options)
}

func (i *Infrastructure) setOpts(opts *option.Options) {
	i.opts.Store(opts)
}

//初始zap.logger包
func (i *Infrastructure) InitLogger() {
	logLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	logEncoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		StacktraceKey:  "StacktraceKey",
		LineEnding:     zapcore.DefaultLineEnding, //add "\n" in line end automatically
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeTime:     logEncodeTime,             //show style of time
		EncodeCaller:   zapcore.FullCallerEncoder, //caller info
		EncodeName:     zapcore.FullNameEncoder,
	}
	logConsoleConfig := zap.Config{
		Level:            logLevel,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    logEncoderConfig,
		OutputPaths:      []string{"stdout", i.GetOpts().ProDir + string(os.PathSeparator) + i.GetOpts().ConsoleOutPutPath},
		ErrorOutputPaths: []string{"stderr", i.GetOpts().ProDir + string(os.PathSeparator) + i.GetOpts().ErrOutPutPath},
	}

	coreOption := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		jsonOutputPath := []string{i.GetOpts().ProDir + string(os.PathSeparator) + i.GetOpts().JsonOutPutPath}
		jsonFileWriteSyncer, _, err := zap.Open(jsonOutputPath...)
		if err != nil {
			panic(err)
		}

		DbLoggerWriteSyncer := dbloggercore.NewDbLoggerWriteSyncer(
			i.DB, i.GetOpts().LoggerBufferCap, i.GetOpts().LoggerBufDuration)

		logEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		c := zapcore.NewCore(
			zapcore.NewJSONEncoder(logEncoderConfig),
			zapcore.NewMultiWriteSyncer(DbLoggerWriteSyncer, jsonFileWriteSyncer),
			logLevel,
		)
		return zapcore.NewTee(core, c)
	})

	var err error
	i.Logger, err = logConsoleConfig.Build(coreOption)
	if err != nil {
		panic(err)
	}
	fmt.Println("初始化zap.logger完成。")
}

func logEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
