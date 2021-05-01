package dag

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestDagRoot(t *testing.T) {
	root := NewTask("root", func(ctx context.Context, helper TaskHelper) (interface{}, error) {
		fmt.Println("root")
		return nil, nil
	})
	dg := NewDag(root)
	dg.Do(context.Background(), false)
}

func TestDagPool(t *testing.T) {
	Config().SetPool()
	locker := sync.Mutex{}
	cntMap := make(map[interface{}]int)
	wg := sync.WaitGroup{}
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				kNum := rand.Intn(100)
				rt := NewTask(fmt.Sprintf("test-%d-%d", index, j), func(ctx context.Context, helper TaskHelper) (interface{}, error) {
					time.Sleep(time.Microsecond * time.Duration(rand.Intn(51)))
					return nil, nil
				})
				cnt := 1
				locker.Lock()
				if val, ok := cntMap[rt]; ok {
					cnt = val + 1
				}
				cntMap[rt] = cnt
				locker.Unlock()
				for k := 0; k < kNum; k++ {
					subTk := NewTask(fmt.Sprintf("sub-%d-%d", index, k), func(ctx context.Context, helper TaskHelper) (interface{}, error) {
						time.Sleep(time.Microsecond * 5)
						return nil, nil
					})

					locker.Lock()
					cnt := 2
					if val, ok := cntMap[subTk]; ok {
						cnt = val + 2
					}
					cntMap[rt] = cnt
					locker.Unlock()
				}
				da := NewDag(rt)
				da.Do(context.Background(), false)
				da.Destory()
			}
		}(i)
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)*10))
	}
	wg.Wait()
	fmt.Println(cntMap)
}

func TestDagTasks_AddRootTasks(t *testing.T) {
	root := NewTask("root", func(ctx context.Context, helper TaskHelper) (interface{}, error) {
		ret, err := helper.GetSubTaskRet("sub-1")
		if err != nil {
			return nil, err
		}
		arr1 := ret.([]int)
		ret, err = helper.GetSubTaskRet("sub-2")
		if err != nil {
			return nil, err
		}
		arr2 := ret.([]int)
		// 做一个归并
		fmt.Println(arr1)
		fmt.Println(arr2)
		fmt.Println("-----------------------------------------------")
		fmt.Println(helper.GetSubTaskRetForIndex(0))
		fmt.Println(helper.GetSubTaskRetForIndex(1))
		fmt.Println(helper.GetSubTaskRetForIndex(2))
		return nil, nil
	})
	num1 := []int{3, 2, 4, 5, 16, 6, 5, 99, 1}
	sub1 := NewTask("sub-1", func(ctx context.Context, helper TaskHelper) (interface{}, error) {
		for i := 0; i < helper.GetSubTaskSize(); i++ {
			res, _ := helper.GetSubTaskRetForIndex(i)
			val := res.([]int)
			num1 = append(num1, val...)
		}
		sort.Ints(num1)
		return num1, nil
	})

	sub11 := NewTask("sub-1-1", func(ctx context.Context, helper TaskHelper) (interface{}, error) {
		return []int{1, 3, 5, 7}, nil
	})

	sub12 := NewTask("sub-1-2", func(ctx context.Context, helper TaskHelper) (interface{}, error) {

		return []int{6, 4, 2, 8}, nil
	})

	num2 := []int{9, 7, 8, 6, 5, 4, 3, 2, 1}
	sub2 := NewTask("sub-2", func(ctx context.Context, helper TaskHelper) (interface{}, error) {
		sort.Ints(num2)
		return num2, nil
	})
	root.AddSubTask(sub1)
	root.AddSubTask(sub2)
	sub1.AddSubTask(sub11)
	sub1.AddSubTask(sub12)
	//sub1.AddSubTask(root)
	dag := NewDag(root)
	fmt.Println(dag.Do(context.Background(), false))
}

func TestCAS(t *testing.T) {
	tk := &task{}
	fmt.Println(tk.state)
	fmt.Println(tk.casSetStatus(tk.state, TaskFinish))
	fmt.Println(tk.state)
}
