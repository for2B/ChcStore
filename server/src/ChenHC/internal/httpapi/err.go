package httpapi

type Err struct {
	Code int
	Msg  string
}

func (e Err) Error() string {
	return e.Msg
}

func (e Err) GetCode() int {
	return e.Code
}

func NewErr(code int, msg string, err error) error {
	if err != nil {
		err_Err, ok := err.(Err)
		if !ok {
			// 说明是原生err
			return Err{Code: code, Msg: msg + "  ::  " + err.Error()}
		}
		return Err{Code: err_Err.Code, Msg: msg + "  ::  " + err_Err.Error()}
	}
	return Err{Code: code, Msg: msg}
}
