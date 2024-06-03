package http

import (
	"app/dto"
	"app/errors"
	"app/pkg/apperror"
	"app/pkg/logger"
	"app/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
//
//	@Summary	Register
//	@Description
//	@Tags		authentications
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.AccountRegisterReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/v1/auth/register [post]
func (s *Server) Register(ctx *gin.Context) {
	req, err := dto.AccountRegisterReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.AuthSvc.Register(ctxFromGin(ctx), req)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, gin.H{"status": "OK"})
}

// Login godoc
//
//	@Summary	Login
//	@Description
//	@Tags		authentications
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.AccountLoginReq	true	"request"
//	@Success	200		{object}	dto.AccessToken
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/v1/auth/login [post]
func (s *Server) Login(ctx *gin.Context) {
	req, err := dto.AccountLoginReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	accessToken, err := s.AuthSvc.Login(ctxFromGin(ctx), req)
	if errors.AccountNotFound().Is(err) {
		logger.For(ctxFromGin(ctx)).Error(apperror.As(err).StackTrace())
		abortWithStatusError(ctx, 400, errors.UserNotRegister())
		return
	}
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, dto.AccessToken{
		AccessToken: accessToken,
	})
}

// Logout godoc
//
//	@Summary	Logout
//	@Description
//	@Tags		authentications
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Success	200
//	@Failure	400	{object}	dto.HTTPResp
//	@Failure	500	{object}	dto.HTTPResp
//	@Router		/api/v1/auth/logout [post]
func (s *Server) Logout(ctx *gin.Context) {
	bearerToken, ok := utils.GetBearerAuth(ctx)
	if !ok {
		abortWithStatusError(ctx, 400, apperror.NewError(errors.CodeTokenError, "empty token"))
		return
	}
	err := s.AuthSvc.Logout(ctxFromGin(ctx), bearerToken)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}
