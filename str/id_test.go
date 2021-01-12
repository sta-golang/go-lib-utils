package str

import (
	"fmt"
	"github.com/rs/xid"
	"testing"
)

func TestXID(t *testing.T) {
	fmt.Println(XID())
	fmt.Println(XID())
	fmt.Println(XID())
	id := xid.New()
	fmt.Println(id.String())
	fmt.Println(id.String())
}
