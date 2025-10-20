package errs

type ErrorResp struct {
	Code int
	Msg  string
	Err  error
}

func NewError(code int, msg string, err error) *ErrorResp {
	return &ErrorResp{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}
