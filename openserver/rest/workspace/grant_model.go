package workspace

import (
	"common"
	"openserver/model"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type GrantModelHandler struct {
	rest.Handler[GrantModelRequest]
}

type GrantModelRequest struct {
	model.UsageLimit `json:",inline"`
}

func NewGrantModelHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &GrantModelHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *GrantModelHandler) Handle() {
	req := &h.Request
	ctx := h.GetContext()
	userId := h.GetFromUser()

	// 工作空间是否属于该用户
	workspace, err := service.Workspace().FindByID(ctx, req.WorkspaceID)
	if err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}

	if workspace.UserID != userId {
		h.SetError(common.WorkspaceNotFound, "workspace owner error")
		return
	}

	if err := service.Workspace().GrantModel(ctx, &req.UsageLimit); err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}
}
