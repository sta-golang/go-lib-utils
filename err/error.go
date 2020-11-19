package err

import "fmt"

type Error struct {
	Code int   `json:"code"`
	Err  error `json:"err"`
}

func (e *Error) getKey() string {
	return fmt.Sprintf("%d-%s", e.Code, e.Err)
}
