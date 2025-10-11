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
	userId := h.GetFromUser()
	err := service.User().Create(ctx, userId, req.NickName, req.RequestLimit, req.TokenLimit)
	if err != nil {
		h.SetError(common.GetErrorCode(err, common.UserCreateError), err.Error())
		return
	}
}

func NewCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &CreateHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}
