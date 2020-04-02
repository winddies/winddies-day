package code

type Code int

const (
	OK          Code = 200
	ResultError Code = 400
	ServerError Code = 500

	LoginFailed   Code = 5000
	RegisterError Code = 5001
)

var errCodeMap = map[Code]string{
	LoginFailed:   "登陆失败，密码或者账号验证错误",
	RegisterError: "注册失败, 账号已经存在",
	ServerError:   "服务器错误",
}

func (code Code) Int() int {
	return int(code)
}

func (code Code) Error() string {
	return errCodeMap[code]
}
