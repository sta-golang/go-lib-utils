package log

import (
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkFileLogAs(b *testing.B) {
	//now := time.Now()
	log := NewFileLogAndAsync(DefaultFileLogConfigForAloneWriter(
		[]string{GetLevelName(INFO), GetLevelName(WARNING), GetLevelName(ERROR)}), time.Second*3)
	for i := 0; i < b.N; i++ {
		log.Info("hello", "world", "golang")
		log.Warn("hello", "world", "golang")
		log.Error("hello", "world", "golang")

	}
	//fmt.Println(time.Now().Sub(now).Milliseconds())
}

func BenchmarkNewFileLogReload(b *testing.B) {
	b.ResetTimer()
	w := writerHelper{
		level:      "all",
		openFile:   atomic.Value{},
		openDate:   "",
		openTime:   0,
		writerSize: 0,
		lock:       sync.Mutex{},
		target: &fileLogWriter{
			helpers:  nil,
			fileDir:  "./sta",
			closeCh:  make(chan *os.File, 1024),
			dayAge:   7,
			fileName: "sta",
			maxSize:  0,
			initFlag: false,
		},
	}
	for i := 0; i < b.N; i++ {
		w.doReLoadFile()
	}
}

func BenchmarkFileLogSy(b *testing.B) {
	log := NewFileLog(DefaultFileLogConfigForAloneWriter(
		[]string{GetLevelName(INFO), GetLevelName(ERROR)}))
	for i := 0; i < b.N; i++ {
		log.Info("hello", "world", "golang")
		log.Warn("hello", "world", "golang")
		log.Error("hello", "world", "golang")
	}
}

func BenchmarkNewFileLogParallel(b *testing.B) {
	log := NewFileLog(DefaultFileLogConfigForAloneWriter(
		[]string{GetLevelName(INFO), GetLevelName(ERROR)}))
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				log.Info("hello", "world", "golang")
			} else {
				log.Error("hello", "world", "golang")
			}
			i++
		}
	})
}

func BenchmarkConsoleLog(b *testing.B) {
	log := NewConsoleLog(INFO, "test")
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			log.Info("hello", "world", "golang")
		} else {
			log.Error("hello", "world", "golang")
		}
	}
}
