package err

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestErrBufferAsync(t *testing.T) {
	InitErrBuffer(true)
	start := sync.WaitGroup{}
	end := sync.WaitGroup{}
	start.Add(1)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		end.Add(1)
		msg := "testErr"
		wg.Add(1)
		go func(code int) {
			defer wg.Done()
			testI := code
			code = code % 25
			end.Done()
			start.Wait()
			PutErr(Error{
				Code: code,
				Err:  errors.New(msg),
			}, "hello - "+string(testI))
		}(i)
	}
	end.Wait()
	start.Done()
	wg.Wait()
	strs := Clean()

	for _, str := range strs {
		fmt.Println(str)
	}
	fmt.Println(len(strs))
}
