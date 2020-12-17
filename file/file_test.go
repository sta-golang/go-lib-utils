package file

import "testing"

func BenchmarkPathExists(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PathExists("file.go")
	}
}
