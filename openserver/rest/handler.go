package rest

import (
	"common"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 请求处理接口
type HandlerInterface interface {
	OnRequest(contex *gin.Context)
}

// 任务处理接口
type TaskInterface interface {
	Handle()
}

type Handler[T any] struct {
	Task       TaskInterface   // 任务处理接口
	Context    *gin.Context    // HTTP上下文
	Request    T               // 请求体
	Response   common.Response // 响应体
	StatusCode int             // 响应状态字
}

func (h *Handler[T]) OnRequest(context *gin.Context) {

	h.Context = context
	h.StatusCode = http.StatusOK
	h.Response.Code = common.Success

	// 解析请求数据
	if err := context.ShouldBind(&h.Request); err != nil {
		h.SetError(common.RequestDataError, err.Error())
		h.SendResponse()
		return
	}

	// 处理请求
	if h.Task != nil {
		defer h.checkPanic()
		h.Task.Handle()
	}

	// 响应请求
	h.SendResponse()
}

func (h *Handler[T]) SendResponse() {
	h.Context.JSON(h.StatusCode, h.Response)
}

func (h *Handler[T]) SetTaskHandler(handler TaskInterface) {
	h.Task = handler
}

func (h *Handler[T]) SetError(code int, msg string) {
	h.Response.Code = code
	h.Response.Msg = msg
}

func (h *Handler[T]) SetStatusCode(statusCode int) {
	h.StatusCode = statusCode
}

func (h *Handler[T]) SetResponseData(data any) {
	h.Response.Data = data
}

func (h *Handler[T]) GetFromUser() string {
	return h.Context.GetString("fromUser")
}

//lint:ignore U1000 Ignore unused function
func (h *Handler[T]) hasJsonBody() bool {
	request := h.Context.Request
	if request.Body == nil || request.Body == http.NoBody {
		return false
	}

	contentType := h.Context.GetHeader("Content-Type")
	return strings.HasPrefix(strings.ToLower(contentType), "application/json")
}

//lint:ignore U1000 Ignore unused function
func (h *Handler[T]) hasQueryParams() bool {
	return len(h.Context.Request.URL.Query()) > 0
}

func (h *Handler[T]) checkPanic() {

	e := recover()
	if e == nil {
		return
	}

	err, ok := e.(error)
	if !ok {
		h.SetError(common.HandleError, "An unkown exception occurs")
		return
	}

	if commError, ok := e.(*common.Error); ok {
		h.SetError(commError.Code, commError.Msg)
	} else {
		h.SetError(common.HandleError, err.Error())
	}
}
