package http

import (
	"app/dto"

	"github.com/gin-gonic/gin"
)

// GetProfileByPhoneNumber godoc
//
//	@Summary	GetProfileByPhoneNumber
//	@Description
//	@Tags		account
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	entity.Account
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/account/profile/{phoneNumber} [get]
func (s *Server) GetProfileByPhoneNumber(ctx *gin.Context) {
	phoneNumber := ctx.Param("phoneNumber")
	account, err := s.AccountSvc.GetProfileByPhoneNumber(ctxFromGin(ctx), phoneNumber)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, account)
}

// GetProfileByUserID godoc
//
//	@Summary	GetProfileByUserID
//	@Description
//	@Tags		account
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	entity.Account
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/account/profile/userID/{userID} [get]
func (s *Server) GetProfileByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	account, err := s.AccountSvc.GetProfileByUserID(ctxFromGin(ctx), userID)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, account)
}

// GetProfile godoc
//
//	@Summary	GetProfile
//	@Description
//	@Tags		account
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	entity.Account
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/v1/account/info [get]
func (s *Server) GetProfile(ctx *gin.Context) {
	account, err := s.AccountSvc.GetProfile(ctxFromGin(ctx))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, account)
}

// CheckPhoneNumber godoc
//
//	@Summary	CheckPhoneNumber
//	@Description
//	@Tags		account
//	@Produce	json
//	@Param		request	body		dto.AccountCheckPhoneNumberReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/v1/account/check-phone [post]
func (s *Server) CheckPhoneNumber(ctx *gin.Context) {
	req, err := dto.AccountCheckPhoneNumberReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.AccountSvc.CheckPhoneNumber(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// ResetPassword godoc
//
//	@Summary	ResetPassword
//	@Description
//	@Tags		account
//	@Produce	json
//	@Param		request	body		dto.AccountResetPasswordReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/v1/account/reset-password [put]
func (s *Server) ResetPassword(ctx *gin.Context) {
	req, err := dto.AccountResetPasswordReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.AccountSvc.ResetPassword(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// ChangePassword godoc
//
//	@Summary	ChangePassword
//	@Description
//	@Tags		account
//	@Produce	json
//	@Param		request	body		dto.AccountChangePasswordReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/v1/account/change-password [put]
func (s *Server) ChangePassword(ctx *gin.Context) {
	req, err := dto.AccountChangePasswordReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.AccountSvc.ChangePassword(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// ChangeAvatar godoc
//
//	@Summary	ChangeAvatar
//	@Description
//	@Tags		account
//	@Produce	json
//	@Param		request	body		dto.AccountChangeAvatarReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/v1/account/change-avatar [put]
func (s *Server) ChangeAvatar(ctx *gin.Context) {
	req, err := dto.AccountChangeAvatarReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.AccountSvc.ChangeAvatar(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}
