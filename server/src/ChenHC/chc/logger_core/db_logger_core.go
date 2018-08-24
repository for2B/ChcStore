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
	loggerBuffer chan *bytes.Buffer
	sync.Mutex
	dB      *sql.DB
	bufPool sync.Pool
	syncReq chan struct{}
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

	go func() {
		ticker := time.NewTicker(time.Duration(LoggerBufDuration) * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				t.SyncRequest()
			}
		}
	}()

	go func() {
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
	case l.loggerBuffer <- b:
	default:
		l.SyncRequest()
		l.loggerBuffer <- b
	}
	return len(p), nil
}

func (l LoggerCore) Sync() error {
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
