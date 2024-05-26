package http

import (
	"app/dto"

	"github.com/gin-gonic/gin"
)

// CreateChat godoc
//
//	@Summary	CreateChat
//	@Description
//	@Tags		chats
//	@Produce	json
//	@Param		request	body		dto.AccessChatReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/chats [post]
func (s *Server) CreateChat(ctx *gin.Context) {
	req, err := dto.AccessChatReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.ChatSvc.CreateChat(ctxFromGin(ctx), req.ToChat(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// CreateGroup godoc
//
//	@Summary	CreateGroup
//	@Description
//	@Tags		chats
//	@Produce	json
//	@Param		request	body		dto.CreateGroupReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/chats/group [post]
func (s *Server) CreateGroup(ctx *gin.Context) {
	req, err := dto.CreateGroupReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.ChatSvc.CreateGroup(ctxFromGin(ctx), req.ToChat(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// GetChatList godoc
//
//	@Summary	GetChatList
//	@Description
//	@Tags		chats
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		page			query		int		false	"page of paging"
//	@Param		page_size		query		int		false	"size of page of paging"
//	@Param		sort			query		string	false	"name of field need to sort"
//	@Param		sort_type		query		string	false	"sort desc or asc"
//	@Param		search			query		string	false	"keyword to search in model"
//	@Success	200				{object}	dto.ChatListResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/chats [get]
func (s *Server) GetChatList(ctx *gin.Context) {
	params, err := dto.QueryParams{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	err = params.Validate(dto.ChatResp{})
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	chats, total, err := s.ChatSvc.GetChatList(ctxFromGin(ctx), params.ToRepoQueryParams())
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	var list = []*dto.ChatResp{}
	for _, u := range chats {
		list = append(list, dto.ChatResp{}.FromChat(u))
	}
	res := dto.ChatListResp{
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		List:     list,
	}
	ctx.AbortWithStatusJSON(200, res)
}

// RenameGroup godoc
//
//	@Summary	RenameGroup
//	@Description
//	@Tags		chats
//	@Produce	json
//	@Param		request	body		dto.RenameGroupReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/chats/group/rename [put]
func (s *Server) RenameGroup(ctx *gin.Context) {
	req, err := dto.RenameGroupReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.ChatSvc.RenameGroup(ctxFromGin(ctx), req.ToChat(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// AddToGroup godoc
//
//	@Summary	RenameChat
//	@Description
//	@Tags		chats
//	@Produce	json
//	@Param		request	body		dto.AddToGroupReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/chats/groupAdd [put]
func (s *Server) AddToGroup(ctx *gin.Context) {
	req, err := dto.AddToGroupReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.ChatSvc.AddToGroup(ctxFromGin(ctx), req.ToChat(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// RemoveFromGroup godoc
//
//	@Summary	RemoveFromGroup
//	@Description
//	@Tags		chats
//	@Produce	json
//	@Param		request	body		dto.RemoveFromGroupReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/chats/groupRemove [put]
func (s *Server) RemoveFromGroup(ctx *gin.Context) {
	req, err := dto.RemoveFromGroupReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.ChatSvc.RemoveFromGroup(ctxFromGin(ctx), req.ToChat(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}
