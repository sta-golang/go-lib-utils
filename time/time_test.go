package time

import (
	"fmt"
	"testing"
)

func TestTimeUnit(t *testing.T) {
	fmt.Println(GetNowDateStr())
	fmt.Println(GetNowDateTimeStr())
	fmt.Println(GetNowTimeStr())
}
