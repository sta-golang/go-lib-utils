package structs

import (
	"fmt"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	toMap, _ := structs{}.toMap(&User{Name: "123", Password: "123", Leader: Leader{LeaderName: "123"}}, "json")
	fmt.Print(toMap)
}
