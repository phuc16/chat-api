package service

import (
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"app/repository"
	"context"
	"time"
)

type ChatService struct {
	ChatRepo IChatRepo
}

func NewChatService(chatRepo IChatRepo) *ChatService {
	return &ChatService{ChatRepo: chatRepo}
}

func (s *ChatService) CreateChat(ctx context.Context, e *entity.Chat) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	chat := &entity.Chat{
		ID:        utils.NewID(),
		ChatName:  "sender",
		IsGroup:   e.IsGroup,
		Users:     e.Users,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.ChatRepo.SaveChat(ctx, chat)
	return
}

func (s *ChatService) CreateGroup(ctx context.Context, e *entity.Chat) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	if len(e.Users) < 2 {
		return nil, errors.ChatGroupNotEnoughUser()
	}
	chat := &entity.Chat{
		ID:        utils.NewID(),
		ChatName:  "sender",
		IsGroup:   e.IsGroup,
		Users:     e.Users,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.ChatRepo.SaveChat(ctx, chat)
	return
}

func (s *ChatService) GetChatList(ctx context.Context, query *repository.QueryParams) (res []*entity.Chat, total int64, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.ChatRepo.GetChatList(ctx, query)
}

func (s *ChatService) RenameGroup(ctx context.Context, e *entity.Chat) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	dbChat, err := s.ChatRepo.GetChatById(ctx, e.ID)
	if err != nil {
		return
	}
	dbChat.ChatName = e.ChatName
	dbChat.UpdatedAt = time.Now()
	err = s.ChatRepo.UpdateChat(ctx, dbChat)
	return
}

func (s *ChatService) AddToGroup(ctx context.Context, chat *entity.Chat) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	err = s.ChatRepo.AddToGroup(ctx, chat)
	if err != nil {
		return
	}
	return
}

func (s *ChatService) RemoveFromGroup(ctx context.Context, chat *entity.Chat) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	err = s.ChatRepo.RemoveFromGroup(ctx, chat)
	if err != nil {
		return
	}
	return
}
