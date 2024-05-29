package service

// type GroupSocketService struct {
// 	UserRepo      IUserRepo
// 	ChatRepo      IChatRepo
// 	GroupRepo     IGroupRepo
// 	VotingRepo    IVotingRepo
// 	UserSocketSvc IUserSocketSvc
// 	ChatSocketSvc IChatSocketSvc
// }

// func NewGroupSocketService(userRepo IUserRepo, chatRepo IChatRepo, GroupRepo IGroupRepo, votingRepo IVotingRepo, chatSocketSvc IChatSocketSvc, userSocketSvc IUserSocketSvc) *GroupSocketService {
// 	return &GroupSocketService{
// 		UserRepo:      userRepo,
// 		ChatRepo:      chatRepo,
// 		GroupRepo:     GroupRepo,
// 		VotingRepo:    votingRepo,
// 		UserSocketSvc: userSocketSvc,
// 		ChatSocketSvc: chatSocketSvc,
// 	}
// }

// func (s *GroupSocketService) Create(ctx context.Context, arrayID []string, info dto.CreateGroupDTO) error {
// 	log.Printf("*** enter create new chat \n***")
// 	log.Printf("* arrayID: %v info: %v *", arrayID, info)

// 	group := entity.Group{
// 		ID:         utils.NewID(),
// 		ChatName:   info.ChatName,
// 		Owner:      info.Owner,
// 		Members:    info.Members,
// 		CreatedAt:  time.Now(),
// 		ChatAvatar: info.Avatar,
// 		Setting: entity.GroupSetting{
// 			ChangeChatNameAndAvatar: true,
// 			PinMessages:             true,
// 			SendMessages:            true,
// 			MembershipApproval:      true,
// 			CreateNewPolls:          true,
// 		},
// 	}
// 	err := s.GroupRepo.SaveGroup(ctx, &group)
// 	if err != nil {
// 		return err
// 	}
// 	conversation := entity.Conversation{
// 		ID:             utils.NewID(),
// 		ChatID:         group.ID,
// 		ID_UserOrGroup: group.ID,
// 		ChatName:       info.ChatName,
// 		ChatAvatar:     info.Avatar,
// 		Type:           entity.TYPE_GROUP,
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 	}

// 	err = s.UserSocketSvc.AppendConversationToMultiple(ctx, arrayID, conversation)
// 	if err != nil {
// 		return err
// 	}
// 	return s.ChatSocketSvc.Create(ctx, group.ID)
// }

// func (s *GroupSocketService) Delete(ctx context.Context, idChat string) error {
// 	log.Println("*** enter delete group \n***")
// 	log.Printf("* idChat: %v *", idChat)
// 	group, err := s.GroupRepo.FindGroupByID(ctx, idChat)
// 	if err != nil {
// 		return err
// 	}
// 	arrayID := s.getListID(group)
// 	err = s.UserSocketSvc.RemoveConversationFromMultiple(ctx, arrayID, idChat)
// 	if err != nil {
// 		return err
// 	}
// 	err = s.ChatSocketSvc.Delete(ctx, idChat)
// 	if err != nil {
// 		return err
// 	}
// 	return s.GroupRepo.DeleteGroupByID(ctx, idChat)
// }

// func (s *GroupSocketService) AppendMember(ctx context.Context, info AppendMemberGroupDTO) error {
// 	log.Println("*** enter append member group \n***")
// 	log.Printf("* info: %v *", info)
// 	personInfo := entity.PersonInfo{
// 		UserID:     info.UserID,
// 		UserName:   info.UserName,
// 		UserAvatar: info.UserAvatar,
// 	}
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	arrayID := s.getListID(group)
// 	for _, id := range arrayID {
// 		if id == personInfo.UserID {
// 			return errors.New("CONFLICT")
// 		}
// 	}
// 	updatedCount, err := s.GroupRepo.AppendMember(ctx, info.IDChat, personInfo)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("append member err: %w", err)

// 	}
// 	conversation := entity.Conversation{
// 		ID:             utils.NewID(),
// 		ChatID:         info.ID,
// 		ID_UserOrGroup: group.ID,
// 		ChatName:       group.ChatName,
// 		ChatAvatar:     group.ChatAvatar,
// 		Type:           entity.TYPE_GROUP,
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 	}
// 	return s.UserSocketSvc.AppendConversation(ctx, info.UserID, conversation)
// }

// func (s *GroupSocketService) AppendAdmin(ctx context.Context, info AppendMemberGroupDTO) error {
// 	log.Println("*** enter append member group \n***")
// 	log.Printf("* info: %v *", info)
// 	personInfo := entity.PersonInfo{
// 		UserID:     info.UserID,
// 		UserName:   info.UserName,
// 		UserAvatar: info.UserAvatar,
// 	}
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	arrayID := s.getListIDAdmin(group)
// 	for _, id := range arrayID {
// 		if id == personInfo.UserID {
// 			return errors.New("CONFLICT")
// 		}
// 	}
// 	_, err = s.GroupRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
// 		updatedCount, err := s.GroupRepo.AppendAdmin(ctx, info.IDChat, personInfo)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return nil, fmt.Errorf("append admin err: %w", err)
// 		}
// 		updatedCount, err = s.GroupRepo.RemoveMember(ctx, info.IDChat, personInfo.UserID)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return nil, fmt.Errorf("append admin err: %w", err)
// 		}
// 		return
// 	})
// 	return err
// }

// func (s *GroupSocketService) ChangeOwner(ctx context.Context, info AppendMemberGroupDTO) error {
// 	log.Println("*** enter change owner group \n***")
// 	log.Printf("* info: %v *", info)
// 	personInfo := entity.PersonInfo{
// 		UserID:     info.UserID,
// 		UserName:   info.UserName,
// 		UserAvatar: info.UserAvatar,
// 	}
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	owner := group.Owner
// 	_, err = s.GroupRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
// 		updatedCount, err := s.GroupRepo.ChangeOwner(ctx, info.IDChat, personInfo)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return nil, fmt.Errorf("change owner admin err: %w", err)
// 		}
// 		updatedCount, err = s.GroupRepo.AppendMember(ctx, info.IDChat, owner)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return nil, fmt.Errorf("append member err: %w", err)
// 		}
// 		updatedCount, err = s.GroupRepo.RemoveMember(ctx, info.IDChat, info.UserID)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return nil, fmt.Errorf("append member err: %w", err)
// 		}
// 		return
// 	})
// 	return err
// }

// func (s *GroupSocketService) RemoveMember(ctx context.Context, info AppendMemberGroupDTO) error {
// 	log.Println("*** enter remove member group \n***")
// 	log.Printf("* info: %v *", info)
// 	personInfo := entity.PersonInfo{
// 		UserID:     info.UserID,
// 		UserName:   info.UserName,
// 		UserAvatar: info.UserAvatar,
// 	}
// 	_, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	updatedCount, err := s.GroupRepo.RemoveMember(ctx, info.IDChat, personInfo.UserID)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("remove member err: %w", err)
// 	}
// 	return s.UserSocketSvc.RemoveConversation(ctx, personInfo.UserID, info.IDChat)
// }

// func (s *GroupSocketService) RemoveAdmin(ctx context.Context, info AppendMemberGroupDTO) error {
// 	log.Println("*** enter remove admin group \n***")
// 	log.Printf("* info: %v *", info)
// 	personInfo := entity.PersonInfo{
// 		UserID:     info.UserID,
// 		UserName:   info.UserName,
// 		UserAvatar: info.UserAvatar,
// 	}
// 	_, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = s.GroupRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
// 		updatedCount, err := s.GroupRepo.RemoveAdmin(ctx, info.IDChat, personInfo.UserID)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return nil, fmt.Errorf("remove admin err: %w", err)
// 		}
// 		updatedCount, err = s.GroupRepo.AppendMember(ctx, info.IDChat, personInfo)
// 		if err != nil || updatedCount.ModifiedCount == 0 {
// 			return nil, fmt.Errorf("append member err: %w", err)
// 		}
// 		return
// 	})
// 	return err
// }

// func (s *GroupSocketService) UpdateNameChat(ctx context.Context, info ChangeNameChatGroupDTO) error {
// 	log.Println("*** enter update name chat group \n***")
// 	log.Printf("* info: %v *", info)
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	updatedCount, err := s.GroupRepo.UpdateNameChat(ctx, info.IDChat, info.ChatName)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("update name chat err: %w", err)
// 	}
// 	arrID := s.GetListID(group)
// 	return s.UserSocketSvc.UpdateChatNameInConversation(ctx, arrID, info.IDChat, info.ChatName)
// }

// func (s *GroupSocketService) UpdateAvatar(ctx context.Context, info ChangeAvatarGroupDTO) error {
// 	log.Println("*** enter update avatar group \n***")
// 	log.Printf("* info: %v *", info)
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	updatedCount, err := s.GroupRepo.UpdateAvatar(ctx, info.IDChat, info.Avatar)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("update avatar err: %w", err)
// 	}
// 	arrID := s.GetListID(group)
// 	return s.UserSocketSvc.UpdateAvatarInConversation(ctx, arrID, info.IDChat, info.Avatar)
// }

// func (s *GroupSocketService) UpdateSettingChangeChatNameAndAvatar(ctx context.Context, info UpdateSettingGroupDTO) error {
// 	log.Println("*** enter update setting change chat name and avatar group \n***")
// 	log.Printf("* info: %v *", info)
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	group.Setting.ChangeChatNameAndAvatar = info.Value
// 	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, info.IDChat, group.Setting)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("update avatar err: %w", err)
// 	}
// 	return err
// }

// func (s *GroupSocketService) UpdateSettingPinMessages(ctx context.Context, info UpdateSettingGroupDTO) error {
// 	log.Println("*** enter update setting pin messages group \n***")
// 	log.Printf("* info: %v *", info)
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	group.Setting.PinMessages = info.Value
// 	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, info.IDChat, group.Setting)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("update avatar err: %w", err)
// 	}
// 	return err
// }

// func (s *GroupSocketService) UpdateSettingSendMessages(ctx context.Context, info UpdateSettingGroupDTO) error {
// 	log.Println("*** enter update setting send messages group \n***")
// 	log.Printf("* info: %v *", info)
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	group.Setting.SendMessages = info.Value
// 	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, info.IDChat, group.Setting)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("update avatar err: %w", err)
// 	}
// 	return err
// }

// func (s *GroupSocketService) UpdateSettingMembershipApproval(ctx context.Context, info UpdateSettingGroupDTO) error {
// 	log.Println("*** enter update setting membership approval group \n***")
// 	log.Printf("* info: %v *", info)
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	group.Setting.MembershipApproval = info.Value
// 	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, info.IDChat, group.Setting)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("update avatar err: %w", err)
// 	}
// 	return err
// }

// func (s *GroupSocketService) UpdateSettingCreateNewPolls(ctx context.Context, info UpdateSettingGroupDTO) error {
// 	log.Println("*** enter update setting create new polls group \n***")
// 	log.Printf("* info: %v *", info)
// 	group, err := s.GroupRepo.FindGroupByID(ctx, info.IDChat)
// 	if err != nil {
// 		return err
// 	}
// 	group.Setting.CreateNewPolls = info.Value
// 	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, info.IDChat, group.Setting)
// 	if err != nil || updatedCount.ModifiedCount == 0 {
// 		return fmt.Errorf("update avatar err: %w", err)
// 	}
// 	return err
// }

// func (s *GroupSocketService) GetListID(group *entity.Group) []string {
// 	listID := []string{group.Owner.UserID}
// 	for _, member := range group.Members {
// 		listID = append(listID, member.UserID)
// 	}
// 	for _, admin := range group.Admins {
// 		listID = append(listID, admin.UserID)
// 	}
// 	return listID
// }

// func (s *GroupSocketService) GetListIDAdmin(group *entity.Group) []string {
// 	listID := []string{}
// 	for _, admin := range group.Admins {
// 		listID = append(listID, admin.UserID)
// 	}
// 	return listID
// }
