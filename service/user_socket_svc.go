package service

import (
	"app/dto"
	"app/entity"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserSocketService struct {
	UserRepo  IUserRepo
	ChatRepo  IChatRepo
	GroupRepo IGroupRepo
}

func NewUserSocketService(userRepo IUserRepo, chatRepo IChatRepo, groupRepo IGroupRepo) *UserSocketService {
	return &UserSocketService{
		UserRepo:  userRepo,
		ChatRepo:  chatRepo,
		GroupRepo: groupRepo,
	}
}

func (s *UserSocketService) AppendFriendRequests(ctx context.Context, req dto.FriendRequestAddDTO) error {
	log.Println("*** enter append friend request ***")
	log.Printf("* req: %v *", &req)

	sender := entity.FriendRequest{
		UserID:      req.SenderID,
		UserName:    req.SenderName,
		UserAvatar:  req.SenderAvatar,
		Description: req.Description,
		IsSender:    false,
		SendAt:      req.SendAt,
	}

	receiver := entity.FriendRequest{
		UserID:      req.ReceiverID,
		UserName:    req.ReceiverName,
		UserAvatar:  req.ReceiverAvatar,
		Description: req.Description,
		IsSender:    true,
		SendAt:      req.SendAt,
	}
	_, err := s.UserRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
		updatedCount, err := s.UserRepo.AppendFriendRequest(ctx, receiver.UserID, sender)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return updatedCount, fmt.Errorf("append friend request sender: %w", err)
		}
		updatedCount, err = s.UserRepo.AppendFriendRequest(ctx, sender.UserID, receiver)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return updatedCount, fmt.Errorf("append friend request receiver: %w", err)
		}
		return
	})
	if err != nil {
		return err
	}

	arrID := getArrID(req.SenderID, req.ReceiverID)
	if _, err := s.UserRepo.SearchConversation(ctx, req.SenderID, arrID[0], arrID[1]); err != nil {
		if err == mongo.ErrNoDocuments {
			appendConvDTO := dto.AppendConversationDTO{
				UserMessageDTO: dto.UserMessageDTO{
					ID:  req.ID,
					TUM: req.TUM,
				},
				SenderID:       req.SenderID,
				SenderName:     req.SenderName,
				SenderAvatar:   req.SenderAvatar,
				ReceiverID:     req.ReceiverID,
				ReceiverName:   req.ReceiverName,
				ReceiverAvatar: req.ReceiverAvatar,
			}
			if err := s.AppendConversations(ctx, appendConvDTO, entity.TYPE_GROUP); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return s.UpdateTypeConversation(ctx, req.SenderID, req.ReceiverID, entity.TYPE_REQUESTS, entity.TYPE_REQUESTED)
}

func (s *UserSocketService) RemoveFriendRequests(ctx context.Context, req dto.FriendRequestRemoveDTO) error {
	log.Println("*** enter remove friend request ***")
	log.Printf("* req: %v *", &req)

	_, err := s.UserRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
		updatedCount, err := s.UserRepo.RemoveFriendRequest(ctx, req.SenderID, req.ReceiverID)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return updatedCount, fmt.Errorf("remove friend request sender: %w", err)
		}
		updatedCount, err = s.UserRepo.RemoveFriendRequest(ctx, req.ReceiverID, req.SenderID)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return updatedCount, fmt.Errorf("remove friend request receiver: %w", err)
		}
		return
	})
	if err != nil {
		return err
	}

	arrID := getArrID(req.SenderID, req.ReceiverID)
	if _, err := s.UserRepo.SearchConversation(ctx, req.SenderID, arrID[0], arrID[1]); err == nil {
		return s.UpdateTypeConversation(ctx, req.SenderID, req.ReceiverID, entity.TYPE_STRANGER, entity.TYPE_STRANGER)
	}

	return nil
}

func (s *UserSocketService) AcceptFriendRequests(ctx context.Context, req dto.FriendRequestAcceptDTO) error {
	log.Println("*** enter accept friend request ***")
	log.Printf("* req: %v *", &req)

	friendRequestRemoveDTO := dto.FriendRequestRemoveDTO{
		UserMessageDTO: dto.UserMessageDTO{
			ID:  req.ID,
			TUM: req.TUM,
		},
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
	}
	if err := s.RemoveFriendRequests(ctx, friendRequestRemoveDTO); err != nil {
		return err
	}

	arrID := getArrID(req.SenderID, req.ReceiverID)
	if _, err := s.UserRepo.SearchConversation(ctx, req.SenderID, arrID[0], arrID[1]); err != nil {
		if err == mongo.ErrNoDocuments {
			appendConvDTO := dto.AppendConversationDTO{
				UserMessageDTO: dto.UserMessageDTO{
					ID:  req.ID,
					TUM: req.TUM,
				},
				SenderID:       req.SenderID,
				SenderName:     req.SenderName,
				SenderAvatar:   req.SenderAvatar,
				ReceiverID:     req.ReceiverID,
				ReceiverName:   req.ReceiverName,
				ReceiverAvatar: req.ReceiverAvatar,
			}
			if err := s.AppendConversations(ctx, appendConvDTO, entity.TYPE_FRIEND); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return s.UpdateTypeConversation(ctx, req.SenderID, req.ReceiverID, entity.TYPE_FRIEND, entity.TYPE_FRIEND)
}

func (s *UserSocketService) Unfriend(ctx context.Context, req dto.UnfriendDTO) error {
	log.Println("*** enter unfriend request ***")
	log.Printf("* req: %v *", &req)

	return s.UpdateTypeConversation(ctx, req.SenderID, req.ReceiverID, entity.TYPE_STRANGER, entity.TYPE_STRANGER)
}

func (s *UserSocketService) UpdateTypeConversation(ctx context.Context, senderID, receiverID, typeSender, typeReceiver string) error {
	log.Println("*** update type conversation ***")
	log.Printf("* senderID: %s receiverID: %s typeSender: %s typeReceiver: %s *", senderID, receiverID, typeSender, typeReceiver)

	chatID1 := senderID[:18] + receiverID[18:]
	chatID2 := receiverID[:18] + senderID[18:]

	_, err := s.UserRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
		updatedCount, err := s.UserRepo.UpdateTypeConversation(ctx, senderID, chatID1, chatID2, typeSender)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return updatedCount, fmt.Errorf("update type conversation sender: %w", err)
		}
		updatedCount, err = s.UserRepo.UpdateTypeConversation(ctx, receiverID, chatID1, chatID2, typeReceiver)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return updatedCount, fmt.Errorf("update type conversation receiver: %w", err)
		}
		return
	})

	return err
}

func (s *UserSocketService) AppendConversations(ctx context.Context, req dto.AppendConversationDTO, _type string) error {
	log.Println("*** append conversations ***")
	log.Printf("* req: %v *", &req)

	arrID := getArrID(req.SenderID, req.ReceiverID)
	idChat := arrID[0]

	if _, err := s.UserRepo.SearchConversation(ctx, req.SenderID, arrID[0], arrID[1]); err == nil {
		return nil
	}

	chat := &entity.Chat{
		ID:             idChat,
		Deliveries:     []entity.Delivery{},
		Reads:          []entity.Delivery{},
		ChatActivities: []entity.ChatActivity{},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.ChatRepo.SaveChat(ctx, chat); err != nil {
		return fmt.Errorf("new chat: %w", err)
	}

	conversationOurSender := entity.Conversation{
		ChatID:            idChat,
		IDUserOrGroup:     req.ReceiverID,
		ChatName:          req.ReceiverName,
		ChatAvatar:        req.ReceiverAvatar,
		Deliveries:        []entity.Delivery{},
		Reads:             []entity.Delivery{},
		TopChatActivities: []entity.ChatActivity{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	switch _type {
	case entity.TYPE_FRIEND:
		conversationOurSender.Type = entity.TYPE_FRIEND
	case entity.TYPE_STRANGER:
		conversationOurSender.Type = entity.TYPE_STRANGER
	default:
		conversationOurSender.Type = entity.TYPE_REQUESTS
	}
	updatedCount, err := s.UserRepo.AppendConversation(ctx, req.SenderID, conversationOurSender)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("append conversation sender: %w", err)
	}

	conversationOurReceiver := entity.Conversation{
		ChatID:            idChat,
		IDUserOrGroup:     req.SenderID,
		ChatName:          req.SenderName,
		ChatAvatar:        req.SenderAvatar,
		Deliveries:        []entity.Delivery{},
		Reads:             []entity.Delivery{},
		TopChatActivities: []entity.ChatActivity{},
	}
	switch _type {
	case entity.TYPE_FRIEND:
		conversationOurReceiver.Type = entity.TYPE_FRIEND
	case entity.TYPE_STRANGER:
		conversationOurReceiver.Type = entity.TYPE_STRANGER
	default:
		conversationOurReceiver.Type = entity.TYPE_REQUESTED
	}
	updatedCount, err = s.UserRepo.AppendConversation(ctx, req.ReceiverID, conversationOurReceiver)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("append conversation receiver: %w", err)
	}

	return nil
}

func (s *UserSocketService) AppendConversation(ctx context.Context, userID string, conversation entity.Conversation) error {
	log.Println("*** append conversations ***")
	log.Printf("* userID: %v conversation: %v *", userID, conversation)

	if _, err := s.UserRepo.SearchSingleConversation(ctx, userID, conversation.ChatID); err != nil {
		updatedCount, err := s.UserRepo.AppendConversation(ctx, userID, conversation)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return fmt.Errorf("append conversation: %w", err)
		}
	}

	return nil
}

func (s *UserSocketService) AppendConversationToMultiple(ctx context.Context, userID []string, conversation entity.Conversation) error {
	log.Println("*** append conversations ***")
	log.Printf("* userID: %v conversation: %v *", userID, conversation)

	updatedCount, err := s.UserRepo.AppendConversationToMultiple(ctx, userID, conversation)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("append conversation: %w", err)
	}

	return nil
}

func (s *UserSocketService) RemoveConversation(ctx context.Context, userID string, idChat string) error {
	log.Println("*** remove conversations ***")
	log.Printf("* userID: %s idChat: %s *", userID, idChat)

	updatedCount, err := s.UserRepo.RemoveConversation(ctx, userID, idChat)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("remove conversation: %w", err)
	}

	return nil
}

func (s *UserSocketService) RemoveConversationFromMultiple(ctx context.Context, userID []string, idChat string) error {
	log.Println("*** remove conversations ***")
	log.Printf("* userID: %s idChat: %s *", userID, idChat)

	updatedCount, err := s.UserRepo.RemoveConversationFromMultiple(ctx, userID, idChat)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("remove conversation: %w", err)
	}

	return nil
}

func (s *UserSocketService) UpdateConversations(ctx context.Context, chat entity.Chat) error {
	log.Println("*** update conversations ***")
	log.Printf("* chat: %v *", chat)

	if len(chat.ChatActivities) > 0 {
		lastActivity := chat.ChatActivities[len(chat.ChatActivities)-1]
		updatedCount, err := s.UserRepo.UpdateChatActivity(
			ctx,
			chat.ID,
			lastActivity.Timestamp,
			chat.Deliveries,
			chat.Reads,
			chat.ChatActivities,
		)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return fmt.Errorf("update conversation: %w", err)
		}
	}
	return nil
}

func (s *UserSocketService) UpdateChatNameInConversation(ctx context.Context, arrID []string, chatID, chatName string) error {
	log.Println("*** update chat name conversations ***")
	log.Printf("* arrID: %v chatID: %s chatName: %s *", arrID, chatID, chatName)

	updatedCount, err := s.UserRepo.UpdateChatNameInConversation(ctx, arrID, chatID, chatName)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("update chat name in conversation: %w", err)
	}

	return nil
}

func (s *UserSocketService) UpdateAvatarInConversation(ctx context.Context, arrID []string, chatID, newAvatar string) error {
	log.Println("*** update avatar conversations ***")
	log.Printf("* arrID: %v chatID: %v newAvatar: %v *", arrID, chatID, newAvatar)

	updatedCount, err := s.UserRepo.UpdateAvatarInConversationMultiple(ctx, arrID, chatID, newAvatar)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("update avatar in conversation: %w", err)
	}

	return nil
}

func getArrID(sender, receiver string) []string {
	chatID1 := sender[:18] + receiver[18:]
	chatID2 := receiver[:18] + sender[18:]
	return []string{chatID1, chatID2}
}
