package http

import (
	"app/dto"

	"github.com/gin-gonic/gin"
)

// SendMessage godoc
//
//	@Summary	SendMessage
//	@Description
//	@Tags		messages
//	@Produce	json
//	@Param		request	body		dto.SendMessageReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/messages [post]
func (s *Server) SendMessage(ctx *gin.Context) {
	req, err := dto.SendMessageReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.MessageSvc.SendMessage(ctxFromGin(ctx), req.ToMessage(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// GetMessageListByChatId godoc
//
//	@Summary	GetMessageListByChatId
//	@Description
//	@Tags		messages
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		page			query		int		false	"page of paging"
//	@Param		page_size		query		int		false	"size of page of paging"
//	@Param		sort			query		string	false	"name of field need to sort"
//	@Param		sort_type		query		string	false	"sort desc or asc"
//	@Param		search			query		string	false	"keyword to search in model"
//	@Success	200				{object}	dto.MessageListResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/messages/{chat_id} [get]
func (s *Server) GetMessageListByChatId(ctx *gin.Context) {
	chatId := ctx.Param("chat_id")
	params, err := dto.QueryParams{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	err = params.Validate(dto.MessageResp{})
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	messages, total, err := s.MessageSvc.GetMessageListByChatId(ctxFromGin(ctx), chatId, params.ToRepoQueryParams())
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	var list = []*dto.MessageResp{}
	for _, u := range messages {
		list = append(list, dto.MessageResp{}.FromMessage(u))
	}
	res := dto.MessageListResp{
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		List:     list,
	}
	ctx.AbortWithStatusJSON(200, res)
}
