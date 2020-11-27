package log

import (
	"testing"
	"time"
)

func BenchmarkFileLogAs(b *testing.B) {
	log := NewFileLogAndAsync(DefaultFileLogConfigForAloneWriter(
		[]string{GetLevelName(INFO), GetLevelName(ERROR)}), time.Second*3)
	for i := 0; i < b.N; i++ {
		log.Info("hello", "world", "golang")
		log.Warn("hello", "world", "golang")
		log.Error("hello", "world", "golang")
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
