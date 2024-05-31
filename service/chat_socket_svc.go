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

type ChatSocketService struct {
	UserRepo   IUserRepo
	ChatRepo   IChatRepo
	GroupRepo  IGroupRepo
	VotingRepo IVotingRepo
}

func NewChatSocketService(userRepo IUserRepo, chatRepo IChatRepo, groupRepo IGroupRepo, votingRepo IVotingRepo) *ChatSocketService {
	return &ChatSocketService{
		UserRepo:   userRepo,
		ChatRepo:   chatRepo,
		GroupRepo:  groupRepo,
		VotingRepo: votingRepo,
	}
}

func (s *ChatSocketService) Create(ctx context.Context, chatID string) error {
	log.Printf("*** enter create new chat ***\n")
	log.Printf("* %s *", chatID)
	chat := entity.Chat{
		ID:             chatID,
		Deliveries:     []entity.Delivery{},
		Reads:          []entity.Delivery{},
		ChatActivities: []entity.ChatActivity{},
	}
	err := s.ChatRepo.SaveChat(ctx, &chat)
	if err != nil {
		return err
	}
	return nil
}

func (s *ChatSocketService) Delete(ctx context.Context, chatID string) error {
	log.Printf("*** enter delete chat ***\n")
	log.Printf("* %s *", chatID)
	err := s.ChatRepo.DeleteChatByID(ctx, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ChatSocketService) AppendChat(ctx context.Context, chatID string, req dto.MessageAppendDTO) error {
	log.Printf("*** enter append chat ***\n")
	log.Printf("* chatID: %s req: %+v *", chatID, req)
	messageID := req.ID
	chatActivity := entity.ChatActivity{
		UserID:     req.UserID,
		UserName:   req.UserName,
		ParentID:   req.ParentID,
		Contents:   req.Contents,
		Timestamp:  req.Timestamp,
		MessageID:  messageID,
		Recall:     false,
		UserAvatar: req.UserAvatar,
	}
	updatedCount, err := s.ChatRepo.AppendChatActivityByIDChat(ctx, chatID, chatActivity)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("append chat activity by id err: %v", err)

	}
	return s.ChangeReadChat(ctx, chatID, dto.MessageDeliveryDTO{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: req.TCM,
		},
		UserID:     req.UserID,
		MessageID:  messageID,
		UserAvatar: req.UserAvatar,
		UserName:   req.UserName,
	})
}

func (s *ChatSocketService) ChangeDeliveryChat(ctx context.Context, chatID string, req dto.MessageDeliveryDTO) error {
	log.Printf("*** enter change delivery chat ***\n")
	log.Printf("* chatID: %s req: %+v *", chatID, req)
	_, err := s.ChatRepo.SearchDeliveryByUserID(ctx, chatID, req.UserID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		updatedCount, err := s.ChatRepo.AppendDelivery(ctx, chatID, entity.Delivery{
			UserID:     req.UserID,
			MessageID:  req.MessageID,
			UserAvatar: req.UserAvatar,
			UserName:   req.UserName,
		})
		if err != nil || updatedCount.ModifiedCount == 0 {
			return fmt.Errorf("append delivery err: %v", err)

		}
	} else {
		updatedCount, err := s.ChatRepo.ChangeDelivery(ctx, chatID, req.UserID, req.MessageID)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return fmt.Errorf("change delivery err: %v", err)
		}
	}
	return nil
}

func (s *ChatSocketService) ChangeReadChat(ctx context.Context, chatID string, req dto.MessageDeliveryDTO) error {
	log.Printf("*** enter change read chat ***\n")
	log.Printf("* chatID: %s req: %+v *", chatID, req)
	_, err := s.ChatRepo.SearchReadByUserID(ctx, chatID, req.UserID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		updatedCount, err := s.ChatRepo.AppendRead(ctx, chatID, entity.Delivery{
			UserID:     req.UserID,
			MessageID:  req.MessageID,
			UserAvatar: req.UserAvatar,
			UserName:   req.UserName,
		})
		if err != nil || updatedCount.ModifiedCount == 0 {
			return fmt.Errorf("append read err: %v", err)
		}
		_, err = s.ChatRepo.SearchDeliveryByUserID(ctx, chatID, req.UserID)
		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}
		if err == mongo.ErrNoDocuments {
			updatedCount, err := s.ChatRepo.AppendDelivery(ctx, chatID, entity.Delivery{
				UserID:     req.UserID,
				MessageID:  req.MessageID,
				UserAvatar: req.UserAvatar,
				UserName:   req.UserName,
			})
			if err != nil || updatedCount.ModifiedCount == 0 {
				return fmt.Errorf("append delivery err: %v", err)
			}
		} else {
			updatedCount, err := s.ChatRepo.ChangeDelivery(ctx, chatID, req.UserID, req.MessageID)
			if err != nil || updatedCount.ModifiedCount == 0 {
				return fmt.Errorf("change delivery err: %v", err)
			}
		}
	} else {
		updatedCount, err := s.ChatRepo.ChangeRead(ctx, chatID, req.UserID, req.MessageID)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return fmt.Errorf("change read err: %v", err)
		}
		updatedCount, err = s.ChatRepo.ChangeDelivery(ctx, chatID, req.UserID, req.MessageID)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return fmt.Errorf("change delivery err: %v", err)
		}
	}
	return nil
}

func (s *ChatSocketService) AppendHiddenMessage(ctx context.Context, chatID string, req dto.MessageHiddenDTO) error {
	log.Printf("*** enter append hidden message ***\n")
	log.Printf("* chatID: %s req: %+v *", chatID, req)
	updatedCount, err := s.ChatRepo.AppendHiddenMessage(ctx, chatID, req.UserID, req.MessageID)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("append hidden message err: %v", err)
	}
	return nil
}

func (s *ChatSocketService) RecallMessage(ctx context.Context, chatID string, req dto.MessageHiddenDTO) error {
	log.Printf("*** enter recall message ***\n")
	log.Printf("* chatID: %s req: %+v *", chatID, req)
	updatedCount, err := s.ChatRepo.RecallMessage(ctx, chatID, req.MessageID)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("recall message err: %v", err)
	}
	return nil
}

func (s *ChatSocketService) GetChatTop10(ctx context.Context, chatID string) (*entity.Chat, error) {
	log.Printf("*** enter get top 10 message ***\n")
	log.Printf("* chatID: %s *", chatID)
	chat, err := s.ChatRepo.GetChatTop10(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *ChatSocketService) AppendVoter(ctx context.Context, req dto.AppendVoterDTO, chatID string, obj dto.MessageAppendDTO) error {
	log.Printf("*** enter append voter ***\n")
	log.Printf("* req: %v chatID: %s obj: %v *", req, chatID, obj)
	updatedCount, err := s.VotingRepo.AppendVoter(ctx, req.VotingID, req.Name, req.Voter)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("append voter err: %v", err)

	}
	return s.AppendChat(ctx, chatID, obj)
}

func (s *ChatSocketService) ChangeVoting(ctx context.Context, req dto.ChangeVoterDTO, chatID string, obj dto.MessageAppendDTO) error {
	log.Printf("*** enter change voting ***\n")
	log.Printf("* req: %v chatID: %s obj: %v *", req, chatID, obj)
	updatedCount, err := s.VotingRepo.RemoveVoter(ctx, req.VotingID, req.OldName, req.Voter.UserID)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("change voting err: %v", err)
	}
	updatedCount, err = s.VotingRepo.AppendVoter(ctx, req.VotingID, req.NewName, req.Voter)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("append voter err: %v", err)
	}
	return s.AppendChat(ctx, chatID, obj)
}

func (s *ChatSocketService) LockVoting(ctx context.Context, chatID string, req dto.MessageAppendDTO) error {
	log.Printf("*** enter lock voting ***\n")
	log.Printf("* chatID: %s obj: %v *", chatID, req)
	updatedCount, err := s.VotingRepo.LockVoting(ctx, req.Contents[0].Value, true, time.Now())
	if err != nil || updatedCount.ModifiedCount == 0 {
		return fmt.Errorf("lock voting err: %v", err)

	}
	return s.AppendChat(ctx, chatID, req)
}
