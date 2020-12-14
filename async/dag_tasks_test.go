package async

import (
	"context"
	"fmt"
	"sort"
	"testing"
)

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
		return []int{1,3,5,7}, nil
	})

	sub12 := NewTask("sub-1-2", func(ctx context.Context, helper TaskHelper) (interface{}, error) {
		return []int{6,4,2,8}, nil
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
