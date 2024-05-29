package service

// type ChatSocketService struct {
// 	UserRepo   IUserRepo
// 	ChatRepo   IChatRepo
// 	GroupRepo  IGroupRepo
// 	VotingRepo IVotingRepo
// }

// func NewChatSocketService(userRepo IUserRepo, chatRepo IChatRepo, groupRepo IGroupRepo, votingRepo IVotingRepo) *ChatSocketService {
// 	return &ChatSocketService{
// 		UserRepo:   userRepo,
// 		ChatRepo:   chatRepo,
// 		GroupRepo:  groupRepo,
// 		VotingRepo: votingRepo,
// 	}
// }

// func (s *ChatSocketService) Create(ctx context.Context, chatID string) error {
// 	log.Printf("*** enter create new chat \n***")
// 	log.Printf("* %s *", chatID)
// 	chat := entity.Chat{
// 		ID: utils.NewID(),
// 	}
// 	err := s.ChatRepo.SaveChat(ctx, &chat)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *ChatSocketService) Delete(ctx context.Context, chatID string) error {
// 	log.Printf("*** enter delete chat \n***")
// 	log.Printf("* %s *", chatID)
// 	err := s.ChatRepo.DeleteChatByID(ctx, chatID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *ChatSocketService) AppendChat(ctx context.Context, chatID string, info dto.MessageAppendDTO) error {
// 	log.Printf("*** enter append chat \n***")
// 	log.Printf("* chatID: %s info: %+v *", chatID, info)
// 	messageID := info.ID
// 	chatActivity := entity.ChatActivity{
// 		ID:         utils.NewID(),
// 		UserID:     info.UserID,
// 		ParentID:   info.ParentID,
// 		Contents:   info.Contents,
// 		Timestamp:  info.Timestamp,
// 		MessageID:  messageID,
// 		Hidden:     info.Hidden,
// 		Recall:     info.Recall,
// 		UserAvatar: info.UserAvatar,
// 	}
// 	updatedCount, err := s.ChatRepo.AppendChatActivityByIDChat(ctx, chatID, chatActivity)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("append chat activity by id err: %w", err)

// 	}
// 	return s.ChangeReadChat(ctx, chatID, dto.MessageDeliveryDTO{
// 		UserID:     info.UserID,
// 		MessageID:  messageID,
// 		UserAvatar: info.UserAvatar,
// 		UserName:   info.UserName,
// 	})
// }

// func (s *ChatSocketService) ChangeDeliveryChat(ctx context.Context, chatID string, info dto.MessageDeliveryDTO) error {
// 	log.Printf("*** enter change delivery chat \n***")
// 	log.Printf("* chatID: %s info: %+v *", chatID, info)
// 	exists, err := s.ChatRepo.SearchDeliveryByUserID(ctx, chatID, info.UserID)
// 	if err != nil {
// 		return err
// 	}
// 	if exists == nil {
// 		updatedCount, err := s.ChatRepo.AppendDelivery(ctx, chatID, entity.Delivery{
// 			UserID:     info.UserID,
// 			MessageID:  info.MessageID,
// 			UserAvatar: info.UserAvatar,
// 			UserName:   info.UserName,
// 		})
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return fmt.Errorf("append delivery err: %w", err)

// 		}
// 	} else {
// 		updatedCount, err := s.ChatRepo.ChangeDelivery(ctx, chatID, info.UserID, info.MessageID)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return fmt.Errorf("change delivery err: %w", err)
// 		}
// 	}
// 	return nil
// }

// func (s *ChatSocketService) ChangeReadChat(ctx context.Context, chatID string, info dto.MessageDeliveryDTO) error {
// 	log.Printf("*** enter change read chat \n***")
// 	log.Printf("* chatID: %s info: %+v *", chatID, info)
// 	exists, err := s.ChatRepo.SearchReadByUserID(ctx, chatID, info.UserID)
// 	if err != nil {
// 		return err
// 	}
// 	if exists == nil {
// 		updatedCount, err := s.ChatRepo.AppendRead(ctx, chatID, entity.Delivery{
// 			UserID:     info.UserID,
// 			MessageID:  info.MessageID,
// 			UserAvatar: info.UserAvatar,
// 			UserName:   info.UserName,
// 		})
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return fmt.Errorf("append read err: %w", err)
// 		}
// 		exists, err = s.ChatRepo.SearchDeliveryByUserID(ctx, chatID, info.UserID)
// 		if err != nil {
// 			return err
// 		}
// 		if exists == nil {
// 			updatedCount, err := s.ChatRepo.AppendDelivery(ctx, chatID, entity.Delivery{
// 				UserID:     info.UserID,
// 				MessageID:  info.MessageID,
// 				UserAvatar: info.UserAvatar,
// 				UserName:   info.UserName,
// 			})
// 			if err != nil || updatedCount.ModifiedCount == 0 {
// 				return fmt.Errorf("append delivery err: %w", err)
// 			}
// 		} else {
// 			updatedCount, err := s.ChatRepo.ChangeDelivery(ctx, chatID, info.UserID, info.MessageID)
// 			if err != nil || updatedCount.ModifiedCount == 0 {
// 				return fmt.Errorf("change delivery err: %w", err)
// 			}
// 		}
// 	} else {
// 		updatedCount, err := s.ChatRepo.ChangeRead(ctx, chatID, info.UserID, info.MessageID)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return fmt.Errorf("change read err: %w", err)
// 		}
// 		updatedCount, err = s.ChatRepo.ChangeDelivery(ctx, chatID, info.UserID, info.MessageID)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return fmt.Errorf("change delivery err: %w", err)
// 		}
// 	}
// 	return nil
// }

// func (s *ChatSocketService) AppendHiddenMessage(ctx context.Context, chatID string, info dto.MessageHiddenDTO) error {
// 	log.Printf("*** enter append hidden message \n***")
// 	log.Printf("* chatID: %s info: %+v *", chatID, info)
// 	updatedCount, err := s.ChatRepo.AppendHiddenMessage(ctx, chatID, info.UserID, info.MessageID)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("append hidden message err: %w", err)
// 	}
// 	return nil
// }

// func (s *ChatSocketService) RecallMessage(ctx context.Context, chatID string, info dto.MessageHiddenDTO) error {
// 	log.Printf("*** enter recall message \n***")
// 	log.Printf("* chatID: %s info: %+v *", chatID, info)
// 	updatedCount, err := s.ChatRepo.RecallMessage(ctx, chatID, info.MessageID)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("recall message err: %w", err)
// 	}
// 	return nil
// }

// func (s *ChatSocketService) GetChatTop10(ctx context.Context, chatID string) (*entity.Chat, error) {
// 	log.Printf("*** enter get top 10 message \n***")
// 	log.Printf("* chatID: %s *", chatID)
// 	chat, err := s.ChatRepo.GetChatTop10(ctx, chatID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return chat, nil
// }

// func (s *ChatSocketService) AppendVoter(ctx context.Context, info dto.AppendVoterDTO, chatID string, obj dto.MessageAppendDTO) error {
// 	log.Printf("*** enter append voter \n***")
// 	log.Printf("* info: %v chatID: %s obj: %v *", info, chatID, obj)
// 	updatedCount, err := s.VotingRepo.AppendVoter(ctx, info.VotingID, info.Name, info.Voter)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("append voter err: %w", err)

// 	}
// 	return s.AppendChat(ctx, chatID, obj)
// }

// func (s *ChatSocketService) ChangeVoting(ctx context.Context, info dto.ChangeVoterDTO, chatID string, obj dto.MessageAppendDTO) error {
// 	log.Printf("*** enter change voting \n***")
// 	log.Printf("* info: %v chatID: %s obj: %v *", info, chatID, obj)
// 	updatedCount, err := s.VotingRepo.RemoveVoter(ctx, info.VotingID, info.OldName, info.Voter.UserID)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("change voting err: %w", err)

// 	}
// 	updatedCount, err = s.VotingRepo.AppendVoter(ctx, info.VotingID, info.NewName, info.Voter)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("append voter err: %w", err)

// 	}
// 	return s.AppendChat(ctx, chatID, obj)
// }

// func (s *ChatSocketService) LockVoting(ctx context.Context, chatID string, obj dto.MessageAppendDTO) error {
// 	log.Printf("*** enter lock voting \n***")
// 	log.Printf("* chatID: %s obj: %v *", chatID, obj)
// 	votingID, err := primitive.ObjectIDFromHex(obj.Contents[0].Value)
// 	if err != nil {
// 		return err
// 	}
// 	updatedCount, err := s.VotingRepo.LockVoting(ctx, votingID, true, time.Now())
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("lock voting err: %w", err)

// 	}
// 	return s.AppendChat(ctx, chatID, obj)
// }
