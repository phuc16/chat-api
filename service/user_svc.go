package service

import (
	"app/dto"
	"app/entity"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
)

type UserService struct {
	UserRepo       IUserRepo
	TokenRepo      ITokenRepo
	UpdateAsyncSvc IUpdateAsyncSvc
}

func NewUserService(userRepo IUserRepo, tokenRepo ITokenRepo, updateAsyncSvc IUpdateAsyncSvc) *UserService {
	return &UserService{UserRepo: userRepo, TokenRepo: tokenRepo, UpdateAsyncSvc: updateAsyncSvc}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	user := &entity.User{
		ID:             req.ID,
		FriendRequests: []entity.FriendRequest{},
		Conversations:  []entity.Conversation{},
	}
	err = s.UserRepo.SaveUser(ctx, user)
	return
}

func (s *UserService) GetUser(ctx context.Context, id string) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.UserRepo.FindUserByID(ctx, id)
}

func (s *UserService) UpdateAvatarAsync(ctx context.Context, req *dto.UpdateAvatarAsyncReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	s.UpdateAsyncSvc.UpdateAvatarAsync(ctx, req.OldAvatar, req.NewAvatar)
	return
}
