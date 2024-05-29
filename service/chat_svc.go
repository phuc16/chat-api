package service

import (
	"app/dto"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
)

type ChatService struct {
	UserRepo IUserRepo
	ChatRepo IChatRepo
}

func NewChatService(userRepo IUserRepo, chatRepo IChatRepo) *ChatService {
	return &ChatService{UserRepo: userRepo, ChatRepo: chatRepo}
}

func (s *ChatService) CreateChat(ctx context.Context, req *dto.CreateChatReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	chat := entity.Chat{
		ID: req.ID,
	}
	err = s.ChatRepo.SaveChat(ctx, &chat)
	return
}

func (s *ChatService) GetChatActivityFromNToM(ctx context.Context, req *dto.GetChatActivityFromNToMReq) (res []entity.ChatActivity, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.ChatRepo.GetChatActivityFromNToM(ctx, req.ID, req.X, req.Y)
}

func (s *ChatService) SearchByKeyWord(ctx context.Context, req *dto.SearchByKeyWordReq) (res []entity.ChatActivity, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.ChatRepo.SearchByKeyWord(ctx, req.ChatID, req.Key)
}

func (s *ChatService) GetSearch(ctx context.Context, req *dto.GetSearchReq) (res []entity.ChatActivity, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	chatIndexes, err := s.ChatRepo.GetIndexOfMessageID(ctx, req.ChatID, req.MessageID)
	if err != nil {
		return
	}
	if len(chatIndexes) == 0 {
		return nil, errors.ChatNotFound()
	}

	getChatActivityFromNToMReq := &dto.GetChatActivityFromNToMReq{
		ID: req.ChatID,
		X:  0,
		Y:  chatIndexes[0].Index + 11,
	}
	return s.GetChatActivityFromNToM(ctx, getChatActivityFromNToMReq)
}
