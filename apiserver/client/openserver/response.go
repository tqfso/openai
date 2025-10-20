package openserver

import "encoding/json"

type Response struct {
	Data json.RawMessage `json:"data"` // 延迟解析
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
}

func (r *Response) IsSuccess() bool {
	return r.Code == 0
}
