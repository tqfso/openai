package workspace

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
	Name string `json:"name" binding:"required"`
}

type CreateResponse struct {
	ID string `json:"id"`
}

func (h *CreateHandler) Handle() {
	req := h.Request
	ctx := h.GetContext()
	userId := h.GetFromUser()
	id, err := service.Workspace().Create(ctx, userId, req.Name)
	if err != nil {
		h.SetError(common.GetErrorCode(err, common.Failure), err.Error())
		return
	}

	h.SetResponseData(&CreateResponse{ID: id})
}

func NewCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &CreateHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}
