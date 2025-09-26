package common

const (
	Success           int = 0 // 成功
	Failure           int = 1 // 失败
	InnerServerError  int = 2 // 内部服务错误
	RequestDataError  int = 3 // 请求数据错误
	RequestParamError int = 4 // 请求参数错误
	HandleError       int = 5 // 处理请求异常
	HandlerNotFound   int = 6 // 未找到处理器
	AuthError         int = 7 // 认证失败
	InnerAccessError  int = 8 // 内部访问失败

	UserExistError  int = 1000 // 用户已存在
	CreateUserError int = 1001 // 创建用户失败

)

var messages = map[int]string{
	Success: "success",
	Failure: "failure",
}

func GetCodeMessage(code int) string {
	msg := messages[code]
	if len(msg) == 0 {
		return "Unknown"
	}

	return msg
}
