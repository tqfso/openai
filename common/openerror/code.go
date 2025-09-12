package openerror

type Code int

const (
	Success Code = 0 // 成功
	Failure Code = 1 // 失败

)

var messages = map[Code]string{
	Success: "success",
	Failure: "failure",
}

func GetMessage(code Code) string {
	msg := messages[code]
	if len(msg) == 0 {
		return "Unknown"
	}

	return msg
}
