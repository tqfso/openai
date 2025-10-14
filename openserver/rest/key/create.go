package key

import (
	"common"
	"openserver/rest"
	"openserver/service"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateHandler struct {
	rest.Handler[CreateRequest]
}

type CreateRequest struct {
	WorkspaceID string     `json:"workspaceID" binding:"required"`
	Description string     `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

type CreateResponse struct {
	ID string `json:"id"`
}

func NewCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &CreateHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *CreateHandler) Handle() {
	req := h.Request
	ctx := h.GetContext()
	userId := h.GetFromUser()
	id, err := service.ApiKey().Create(ctx, userId, req.WorkspaceID, req.Description, req.ExpiresAt)
	if err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}

	h.SetResponseData(&CreateResponse{ID: id})
}
