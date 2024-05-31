package http

import (
	"app/dto"

	"github.com/gin-gonic/gin"
)

// CreateGroup godoc
//
//	@Summary	CreateGroup
//	@Description
//	@Tags		group
//	@Produce	json
//	@Param		Authorization	header		string				true	"Bearer token"
//	@Param		request			body		dto.CreateGroupReq	true	"request"
//	@Success	200				{object}	dto.HTTPResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/group/create [post]
func (s *Server) CreateGroup(ctx *gin.Context) {
	req, err := dto.CreateGroupReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.GroupSvc.CreateGroup(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, gin.H{"status": "OK"})
}

// GetGroupInfo godoc
//
//	@Summary	GetGroupInfo
//	@Description
//	@Tags		group
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		groupID			query		string	true	"groupID"
//	@Success	200				{object}	entity.Group
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/group/info [get]
func (s *Server) GetGroupInfo(ctx *gin.Context) {
	req, err := dto.GetGroupInfoReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	group, err := s.GroupSvc.GetGroupInfo(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, group)
}
