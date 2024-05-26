package service

import (
	"app/entity"
	"app/pkg/trace"
	"app/pkg/utils"
	"app/repository"
	"context"
	"time"
)

type MessageService struct {
	OtpSvc      IOtpSvc
	MessageRepo IMessageRepo
	ChatRepo    IChatRepo
}

func NewMessageService(messageRepo IMessageRepo, chatRepo IChatRepo) *MessageService {
	return &MessageService{MessageRepo: messageRepo, ChatRepo: chatRepo}
}

func (s *MessageService) SendMessage(ctx context.Context, e *entity.Message) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	message := &entity.Message{
		ID:        utils.NewID(),
		Sender:    e.Sender,
		Message:   e.Message,
		ChatID:    e.ChatID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	res, err = s.MessageRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
		err = s.MessageRepo.SaveMessage(ctx, message)
		if err != nil {
			return
		}
		dbChat, err := s.ChatRepo.GetChatById(ctx, e.ChatID)
		if err != nil {
			return
		}
		dbChat.LatestMessage = *message
		dbChat.UpdatedAt = time.Now()
		err = s.ChatRepo.UpdateChat(ctx, dbChat)
		if err != nil {
			return
		}
		return
	})
	return
}

func (s *MessageService) GetMessageListByChatId(ctx context.Context, chatId string, query *repository.QueryParams) (res []*entity.Message, total int64, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.MessageRepo.GetMessageListByChatId(ctx, chatId, query)
}
