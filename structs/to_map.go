package structs

import (
	"fmt"
	"reflect"
)

type structs struct {
}

type Leader struct {
	LeaderName string `json:"leaderName"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Leader   Leader `json:"leader"`
}

// toMap
func (s structs) toMap(in interface{}, tag string) (map[string]interface{}, error) {
	//反射获取对象类型
	v := reflect.ValueOf(in)
	//如果对象是指针类型,则转化为指针对应的结构体
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	//如果接收的类型不是结构体类型,则抛出异常
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts struct")
	}
	out := make(map[string]interface{})
	queue := make([]interface{}, 0, 1)
	queue = append(queue, in)
	for len(queue) > 0 {
		//将所有结构体放入queue中进行遍历,并截取将遍历的元素出队
		v := reflect.ValueOf(queue[0])
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		queue = queue[1:]
		t := v.Type()
		//获取当前元素所有的子集合元素
		for i := 0; i < v.NumField(); i++ {
			vi := v.Field(i)
			//如果当前元素也是结构体类型,则将元素入队,进行下一轮循环
			if vi.Kind() == reflect.Struct {
				queue = append(queue, vi.Interface())
				continue
			}
			//如果当前元素是指针类型,则获取到指针的值并进行判断
			if vi.Kind() == reflect.Ptr {
				vi = vi.Elem()
				if vi.Kind() == reflect.Struct {
					queue = append(queue, vi.Interface())
				} else {
					//将普通字段直接加入到map中
					ti := t.Field(i)
					if tagValue := ti.Tag.Get(tag); tagValue != "" {
						out[tagValue] = vi.Interface()
					}
				}
				continue
			}
			ti := t.Field(i)
			if tagValue := ti.Tag.Get(tag); tagValue != "" {
				out[tagValue] = vi.Interface()
			}
		}
	}
	return out, nil
}
