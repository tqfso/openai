package user

import (
	"common"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type CreateHandler struct {
	rest.Handler[CreateRequest]
}

type CreateRequest struct {
	NickName     string `json:"nickName"`
	RequestLimit int64  `json:"requestLimit"`
	TokenLimit   int64  `json:"tokenLimit"`
}

func (h *CreateHandler) Handle() {

	req := h.Request
	ctx := h.GetContext()
	id := h.GetFromUser()
	svc := service.NewUserService()
	err := svc.CreateUser(ctx, id, req.NickName, req.RequestLimit, req.TokenLimit)
	if err != nil {
		h.SetError(common.GetErrorCode(err, common.CreateUserError), err.Error())
		return
	}
}

func NewCreateHandler() gin.HandlerFunc {
	h := &CreateHandler{}
	h.SetTaskHandler(h)
	return h.OnRequest
}
