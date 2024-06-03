package service

import (
	"app/dto"
	"app/entity"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type GroupSocketService struct {
	UserRepo      IUserRepo
	ChatRepo      IChatRepo
	GroupRepo     IGroupRepo
	VotingRepo    IVotingRepo
	UserSocketSvc IUserSocketSvc
	ChatSocketSvc IChatSocketSvc
}

func NewGroupSocketService(userRepo IUserRepo, chatRepo IChatRepo, GroupRepo IGroupRepo, votingRepo IVotingRepo, chatSocketSvc IChatSocketSvc, userSocketSvc IUserSocketSvc) *GroupSocketService {
	return &GroupSocketService{
		UserRepo:      userRepo,
		ChatRepo:      chatRepo,
		GroupRepo:     GroupRepo,
		VotingRepo:    votingRepo,
		UserSocketSvc: userSocketSvc,
		ChatSocketSvc: chatSocketSvc,
	}
}

func (s *GroupSocketService) Create(ctx context.Context, arrayID []string, req dto.CreateGroupDTO) ([]string, error) {
	log.Printf("*** enter create new chat ***")
	log.Printf("* arrayID: %v req: %v *", arrayID, req)

	group := entity.Group{
		ID:         req.ID,
		ChatName:   req.ChatName,
		Owner:      req.Owner,
		Admins:     []entity.PersonInfo{},
		Members:    req.Members,
		CreatedAt:  time.Now(),
		ChatAvatar: req.Avatar,
		Setting: entity.GroupSetting{
			ChangeChatNameAndAvatar: true,
			PinMessages:             true,
			SendMessages:            true,
			MembershipApproval:      true,
			CreateNewPolls:          true,
		},
	}
	err := s.GroupRepo.SaveGroup(ctx, &group)
	if err != nil {
		return nil, err
	}
	conversation := entity.Conversation{
		ChatID:            req.ID,
		IDUserOrGroup:     req.ID,
		ChatName:          req.ChatName,
		ChatAvatar:        req.Avatar,
		Type:              entity.TYPE_GROUP,
		Deliveries:        []entity.Delivery{},
		Reads:             []entity.Delivery{},
		TopChatActivities: []entity.ChatActivity{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	err = s.UserSocketSvc.AppendConversationToMultiple(ctx, arrayID, conversation)
	if err != nil {
		return nil, err
	}
	err = s.ChatSocketSvc.Create(ctx, group.ID)
	if err != nil {
		return nil, err
	}
	return s.GetListID(&group), nil
}

func (s *GroupSocketService) Delete(ctx context.Context, idChat string) ([]string, error) {
	log.Println("*** enter delete group ***")
	log.Printf("* idChat: %v *", idChat)
	group, err := s.GroupRepo.FindGroupByID(ctx, idChat)
	if err != nil {
		return nil, err
	}
	arrayID := s.GetListID(group)
	err = s.UserSocketSvc.RemoveConversationFromMultiple(ctx, arrayID, idChat)
	if err != nil {
		return nil, err
	}
	err = s.ChatSocketSvc.Delete(ctx, idChat)
	if err != nil {
		return nil, err
	}
	err = s.GroupRepo.DeleteGroupByID(ctx, idChat)
	if err != nil {
		return nil, err
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) AppendMember(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error) {
	log.Println("*** enter append member group ***")
	log.Printf("* req: %v *", req)
	personInfo := entity.PersonInfo{
		UserID:     req.UserID,
		UserName:   req.UserName,
		UserAvatar: req.UserAvatar,
	}
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	arrayID := s.GetListID(group)
	for _, id := range arrayID {
		if id == personInfo.UserID {
			return nil, errors.New("CONFLICT")
		}
	}
	updatedCount, err := s.GroupRepo.AppendMember(ctx, req.IDChat, personInfo)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return nil, fmt.Errorf("append member err: %w", err)

	}
	conversation := entity.Conversation{
		ChatID:            req.ID,
		IDUserOrGroup:     group.ID,
		ChatName:          group.ChatName,
		ChatAvatar:        group.ChatAvatar,
		Type:              entity.TYPE_GROUP,
		Deliveries:        []entity.Delivery{},
		Reads:             []entity.Delivery{},
		TopChatActivities: []entity.ChatActivity{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = s.UserSocketSvc.AppendConversation(ctx, req.UserID, conversation)
	if err != nil {
		return nil, err
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) AppendAdmin(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error) {
	log.Println("*** enter append member group ***")
	log.Printf("* req: %v *", req)
	personInfo := entity.PersonInfo{
		UserID:     req.UserID,
		UserName:   req.UserName,
		UserAvatar: req.UserAvatar,
	}
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	arrayID := s.GetListIDAdmin(group)
	for _, id := range arrayID {
		if id == personInfo.UserID {
			return nil, errors.New("CONFLICT")
		}
	}
	_, err = s.GroupRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
		updatedCount, err := s.GroupRepo.AppendAdmin(ctx, req.IDChat, personInfo)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return nil, fmt.Errorf("append admin err: %w", err)
		}
		updatedCount, err = s.GroupRepo.RemoveMember(ctx, req.IDChat, personInfo.UserID)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return nil, fmt.Errorf("append admin err: %w", err)
		}
		return
	})
	if err != nil {
		return nil, err
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) ChangeOwner(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error) {
	log.Println("*** enter change owner group ***")
	log.Printf("* req: %v *", req)
	personInfo := entity.PersonInfo{
		UserID:     req.UserID,
		UserName:   req.UserName,
		UserAvatar: req.UserAvatar,
	}
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	owner := group.Owner
	_, err = s.GroupRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
		updatedCount, err := s.GroupRepo.ChangeOwner(ctx, req.IDChat, personInfo)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return nil, fmt.Errorf("change owner admin err: %w", err)
		}
		updatedCount, err = s.GroupRepo.AppendMember(ctx, req.IDChat, owner)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return nil, fmt.Errorf("append member err: %w", err)
		}
		updatedCount, err = s.GroupRepo.RemoveMember(ctx, req.IDChat, req.UserID)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return nil, fmt.Errorf("append member err: %w", err)
		}
		return
	})
	if err != nil {
		return nil, err
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) RemoveMember(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error) {
	log.Println("*** enter remove member group ***")
	log.Printf("* req: %v *", req)
	personInfo := entity.PersonInfo{
		UserID:     req.UserID,
		UserName:   req.UserName,
		UserAvatar: req.UserAvatar,
	}
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	updatedCount, err := s.GroupRepo.RemoveMember(ctx, req.IDChat, personInfo.UserID)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return nil, fmt.Errorf("remove member err: %w", err)
	}
	err = s.UserSocketSvc.RemoveConversation(ctx, personInfo.UserID, req.IDChat)
	if err != nil {
		return nil, err
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) RemoveAdmin(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error) {
	log.Println("*** enter remove admin group ***")
	log.Printf("* req: %v *", req)
	personInfo := entity.PersonInfo{
		UserID:     req.UserID,
		UserName:   req.UserName,
		UserAvatar: req.UserAvatar,
	}
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	_, err = s.GroupRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
		updatedCount, err := s.GroupRepo.RemoveAdmin(ctx, req.IDChat, personInfo.UserID)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return nil, fmt.Errorf("remove admin err: %w", err)
		}
		updatedCount, err = s.GroupRepo.AppendMember(ctx, req.IDChat, personInfo)
		if err != nil || updatedCount.ModifiedCount == 0 {
			return nil, fmt.Errorf("append member err: %w", err)
		}
		return
	})
	if err != nil {
		return nil, err
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) UpdateNameChat(ctx context.Context, req dto.ChangeNameChatGroupDTO) ([]string, error) {
	log.Println("*** enter update name chat group ***")
	log.Printf("* req: %v *", req)
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	updatedCount, err := s.GroupRepo.UpdateNameChat(ctx, req.IDChat, req.ChatName)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return nil, fmt.Errorf("update name chat err: %w", err)
	}
	arrID := s.GetListID(group)
	err = s.UserSocketSvc.UpdateChatNameInConversation(ctx, arrID, req.IDChat, req.ChatName)
	if err != nil {
		return nil, err
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) UpdateAvatar(ctx context.Context, req dto.ChangeAvatarGroupDTO) ([]string, error) {
	log.Println("*** enter update avatar group ***")
	log.Printf("* req: %v *", req)
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	updatedCount, err := s.GroupRepo.UpdateAvatar(ctx, req.IDChat, req.Avatar)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return nil, fmt.Errorf("update avatar err: %w", err)
	}
	arrID := s.GetListID(group)
	err = s.UserSocketSvc.UpdateAvatarInConversation(ctx, arrID, req.IDChat, req.Avatar)
	if err != nil {
		return nil, err
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) UpdateSettingChangeChatNameAndAvatar(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error) {
	log.Println("*** enter update setting change chat name and avatar group ***")
	log.Printf("* req: %v *", req)
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	group.Setting.ChangeChatNameAndAvatar = req.Value
	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, req.IDChat, group.Setting)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return nil, fmt.Errorf("update avatar err: %w", err)
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) UpdateSettingPinMessages(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error) {
	log.Println("*** enter update setting pin messages group ***")
	log.Printf("* req: %v *", req)
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	group.Setting.PinMessages = req.Value
	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, req.IDChat, group.Setting)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return nil, fmt.Errorf("update avatar err: %w", err)
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) UpdateSettingSendMessages(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error) {
	log.Println("*** enter update setting send messages group ***")
	log.Printf("* req: %v *", req)
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	group.Setting.SendMessages = req.Value
	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, req.IDChat, group.Setting)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return nil, fmt.Errorf("update avatar err: %w", err)
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) UpdateSettingMembershipApproval(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error) {
	log.Println("*** enter update setting membership approval group ***")
	log.Printf("* req: %v *", req)
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	group.Setting.MembershipApproval = req.Value
	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, req.IDChat, group.Setting)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return nil, fmt.Errorf("update avatar err: %w", err)
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) UpdateSettingCreateNewPolls(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error) {
	log.Println("*** enter update setting create new polls group ***")
	log.Printf("* req: %v *", req)
	group, err := s.GroupRepo.FindGroupByID(ctx, req.IDChat)
	if err != nil {
		return nil, err
	}
	group.Setting.CreateNewPolls = req.Value
	updatedCount, err := s.GroupRepo.UpdateSetting(ctx, req.IDChat, group.Setting)
	if err != nil || updatedCount.ModifiedCount == 0 {
		return nil, fmt.Errorf("update avatar err: %w", err)
	}
	return s.GetListID(group), nil
}

func (s *GroupSocketService) GetListID(group *entity.Group) []string {
	listID := []string{group.Owner.UserID}
	for _, member := range group.Members {
		listID = append(listID, member.UserID)
	}
	for _, admin := range group.Admins {
		listID = append(listID, admin.UserID)
	}
	return listID
}

func (s *GroupSocketService) GetListIDAdmin(group *entity.Group) []string {
	listID := []string{}
	for _, admin := range group.Admins {
		listID = append(listID, admin.UserID)
	}
	return listID
}
