package dbloggercore

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"sync"
	"time"
)

type LoggerCore struct {
	loggerBuffer chan *bytes.Buffer  //logger缓冲区
	sync.Mutex
	dB      *sql.DB	//数据库
	bufPool sync.Pool //对象池
	syncReq chan struct{}  //传递信号
}

func NewDbLoggerWriteSyncer(db *sql.DB, loggerBufferCap int, LoggerBufDuration int64) *LoggerCore {
	t := &LoggerCore{
		dB:           db,
		loggerBuffer: make(chan *bytes.Buffer, loggerBufferCap),
		bufPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		syncReq: make(chan struct{}),
	}

	go func() {  //开一个定时线程，每个固定LoggerBufDuration时间，就发出同步信号同步日志信息
		ticker := time.NewTicker(time.Duration(LoggerBufDuration) * time.Millisecond)  //NewTicker返回一个新的Ticker，该Ticker包含一个通道字段，并会每隔时间段d就向该通道发送当时的时间
		for {
			select {
			case <-ticker.C:  //Ticker保管一个通道，并每隔一段时间向其传递"tick"。
				t.SyncRequest()
			}
		}
	}()

	go func() { //接受到同步信号后开始同步
		for {
			select {
			case <-t.syncReq:
				t.Sync()
			}
		}
	}()

	return t
}

func (l LoggerCore) Write(p []byte) (n int, err error) {
	b := l.bufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.Write(p)
	select {
	case l.loggerBuffer <- b:  //写入管道
	default: //如果管道满了则开启同步,同步完了之后再写入
		l.SyncRequest()
		l.loggerBuffer <- b
	}
	return len(p), nil
}

func (l LoggerCore) Sync() error { //同步数据，将logger数据写入到数据库
	l.Lock()
	defer l.Unlock()
	txn, err := l.dB.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare(pq.CopyIn("syslog", "log"))
	if err != nil {
		return err
	}

Exit:
	for {
		select {
		case t := <-l.loggerBuffer:
			_, err := stmt.Exec(t.String())
			l.bufPool.Put(t)
			if err != nil {
				fmt.Println(err.Error())
			}
		default:
			break Exit
		}
	}
	if _, err = stmt.Exec(); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}
	if err = txn.Commit(); err != nil {
		return err
	}
	return nil
}

func (l LoggerCore) SyncRequest() {
	select {
	case l.syncReq <- struct{}{}:
	default:
	}
}
