package http

import (
	"app/dto"

	"github.com/gin-gonic/gin"
)

// CreateChat godoc
//
//	@Summary	CreateChat
//	@Description
//	@Tags		chat
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		id				query		string	true	"id"
//	@Success	200				{object}	dto.HTTPResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/chat/create [post]
func (s *Server) CreateChat(ctx *gin.Context) {
	req, err := dto.CreateChatReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.ChatSvc.CreateChat(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// GetChatActivityFromNToM godoc
//
//	@Summary	GetChatActivityFromNToM
//	@Description
//	@Tags		chat
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		id				query		string	true	"id"
//	@Param		x				query		int		true	"x"
//	@Param		y				query		int		true	"y"
//	@Success	200				{object}	[]entity.ChatActivity
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/chat/x-to-y [get]
func (s *Server) GetChatActivityFromNToM(ctx *gin.Context) {
	req, err := dto.GetChatActivityFromNToMReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	chatActivities, err := s.ChatSvc.GetChatActivityFromNToM(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, chatActivities)
}

// SearchByKeyWord godoc
//
//	@Summary	SearchByKeyWord
//	@Description
//	@Tags		chat
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		chatID			query		string	true	"chatID"
//	@Param		key				query		string	true	"key"
//	@Success	200				{object}	[]entity.ChatActivity
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/chat/search-bkw [get]
func (s *Server) SearchByKeyWord(ctx *gin.Context) {
	req, err := dto.SearchByKeyWordReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	chatActivities, err := s.ChatSvc.SearchByKeyWord(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, chatActivities)
}

// GetSearch godoc
//
//	@Summary	GetSearch
//	@Description
//	@Tags		chat
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		chatID			query		string	true	"chatID"
//	@Param		messageID		query		string	true	"messageID"
//	@Success	200				{object}	[]entity.ChatActivity
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/chat/get-search [get]
func (s *Server) GetSearch(ctx *gin.Context) {
	req, err := dto.GetSearchReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	chatActivities, err := s.ChatSvc.GetSearch(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, chatActivities)
}
