package log

import "testing"

func BenchmarkFileLog(b *testing.B) {
	log := NewFileLog(DefaultFileLogConfigForAloneWriter(
		[]string{LEVEL_FLAGS[INFO], LEVEL_FLAGS[ERROR]}))
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			log.Info("hello", "world", "golang")
		} else {
			log.Error("hello", "world", "golang")
		}
	}
}

func BenchmarkNewFileLogParallel(b *testing.B) {
	log := NewFileLog(DefaultFileLogConfigForAloneWriter(
		[]string{LEVEL_FLAGS[INFO], LEVEL_FLAGS[ERROR]}))
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
