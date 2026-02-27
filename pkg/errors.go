package errors

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

const (
	ErrRateLimitExceededCode = iota + 10000
)

var (
	RateLimitExceededError = NewError(ErrRateLimitExceededCode, "请求频繁")
)
