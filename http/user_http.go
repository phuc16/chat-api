package http

import (
	"app/dto"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
//
//	@Summary	CreateUser
//	@Description
//	@Tags		authentications
//	@Produce	json
//	@Param		id	query		string	true	"id"
//	@Success	200	{object}	dto.HTTPResp
//	@Failure	400	{object}	dto.HTTPResp
//	@Failure	500	{object}	dto.HTTPResp
//	@Router		/api/v1/user/create [post]
func (s *Server) CreateUser(ctx *gin.Context) {
	req, err := dto.CreateUserReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.UserSvc.CreateUser(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// GetUser godoc
//
//	@Summary	GetUser
//	@Description
//	@Tags		user
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	entity.User
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/user/info/{id} [get]
func (s *Server) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := s.UserSvc.GetUser(ctxFromGin(ctx), id)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, user)
}

// UpdateAvatarAsync godoc
//
//	@Summary	UpdateAvatarAsync
//	@Description
//	@Tags		user
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		oldAvatar		query		string	true	"oldAvatar"
//	@Param		newAvatar		query		string	true	"newAvatar"
//	@Success	200				{object}	dto.HTTPResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/user/update-avatar-account [get]
func (s *Server) UpdateAvatarAsync(ctx *gin.Context) {
	req, err := dto.UpdateAvatarAsyncReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.UserSvc.UpdateAvatarAsync(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, gin.H{"status": "OK"})
}

// // SuggestFriend godoc
// //
// //	@Summary	SuggestFriend
// //	@Description
// //	@Tags		users
// //	@Produce	json
// //	@Param		Authorization	header		string	true	"Bearer token"
// //	@Success	200				{object}	[]dto.UserResp
// //	@Failure	400				{object}	dto.HTTPResp
// //	@Failure	500				{object}	dto.HTTPResp
// //	@Router		/api/users/friends/suggest [get]
// func (s *Server) SuggestFriend(ctx *gin.Context) {
// 	users, err := s.UserSvc.SuggestFriend(ctxFromGin(ctx), &entity.User{
// 		ID: entity.GetUserFromContext(ctxFromGin(ctx)).ID,
// 	})
// 	if err != nil {
// 		abortWithStatusError(ctx, 400, err)
// 		return
// 	}
// 	var list = []*dto.UserResp{}
// 	for _, u := range users {
// 		list = append(list, dto.UserResp{}.FromUser(u))
// 	}
// 	ctx.AbortWithStatusJSON(200, list)
// }
