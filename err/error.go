package err

import "fmt"

const (
	// error 规范 9开头为框架 后两位为具体的package 后两位为各自包自己规定的编码
	LogErrCode = 90400
)

type Error struct {
	Code int   `json:"code"`
	Err  error `json:"err"`
}

func (e *Error) getKey() string {
	return fmt.Sprintf("%d-%s", e.Code, e.Err)
}

func NewError(code int, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func (e *Error) String() string {
	return fmt.Sprintf("err code : %v, msg : %v", e.Code, e.Err)
}
