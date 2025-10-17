package workspace

import (
	"common"
	"openserver/rest"
	"openserver/service"

	"github.com/gin-gonic/gin"
)

type CancelModelHandler struct {
	rest.Handler[CancelModelRequest]
}

type CancelModelRequest struct {
	WorkspaceID string `json:"workspaceID" binding:"required"`
	ModelName   string `json:"modelName" binding:"required"`
}

func NewCancelModelHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &CancelModelHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *CancelModelHandler) Handle() {
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

	if err := service.Workspace().CancelModel(ctx, req.WorkspaceID, req.ModelName); err != nil {
		h.SetErrorWithDefaultCode(err, common.Failure)
		return
	}
}
