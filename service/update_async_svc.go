package service

import (
	"context"
	"fmt"
	"log"
)

type UpdateAsyncService struct {
	UserRepo  IUserRepo
	ChatRepo  IChatRepo
	GroupRepo IGroupRepo
}

func NewUpdateAsyncService(UserRepo IUserRepo, ChatRepo IChatRepo, GroupRepo IGroupRepo) *UpdateAsyncService {
	return &UpdateAsyncService{UserRepo: UserRepo, ChatRepo: ChatRepo, GroupRepo: GroupRepo}
}

func (s *UpdateAsyncService) UpdateAvatarAsync(ctx context.Context, oldAvatar, newAvatar string) {
	go func() {
		errChan := make(chan error, 7)

		go func() {
			updatedCount, err := s.UserRepo.UpdateAvatarInFriendRequest(ctx, oldAvatar, newAvatar)
			if err != nil || updatedCount.ModifiedCount == 0 {
				err = fmt.Errorf("UpdateAvatarInFriendRequest %v:", err)
			}
			errChan <- err
		}()
		go func() {
			updatedCount, err := s.UserRepo.UpdateAvatarInConversation(ctx, oldAvatar, newAvatar)
			if err != nil || updatedCount.ModifiedCount == 0 {
				err = fmt.Errorf("UpdateAvatarInConversation %v:", err)
			}
			errChan <- err
		}()
		go func() {
			updatedCount, err := s.ChatRepo.UpdateAvatarInDelivery(ctx, oldAvatar, newAvatar)
			if err != nil || updatedCount.ModifiedCount == 0 {
				err = fmt.Errorf("UpdateAvatarInDelivery %v:", err)
			}
			errChan <- err
		}()
		go func() {
			updatedCount, err := s.ChatRepo.UpdateAvatarInRead(ctx, oldAvatar, newAvatar)
			if err != nil || updatedCount.ModifiedCount == 0 {
				err = fmt.Errorf("UpdateAvatarInRead %v:", err)
			}
			errChan <- err
		}()
		go func() {
			updatedCount, err := s.GroupRepo.UpdateAvatarInOwner(ctx, oldAvatar, newAvatar)
			if err != nil || updatedCount.ModifiedCount == 0 {
				err = fmt.Errorf("UpdateAvatarInOwner %v:", err)
			}
			errChan <- err
		}()
		go func() {
			updatedCount, err := s.GroupRepo.UpdateAvatarInAdmins(ctx, oldAvatar, newAvatar)
			if err != nil || updatedCount.ModifiedCount == 0 {
				err = fmt.Errorf("UpdateAvatarInAdmins %v:", err)
			}
			errChan <- err
		}()
		go func() {
			updatedCount, err := s.GroupRepo.UpdateAvatarInMembers(ctx, oldAvatar, newAvatar)
			if err != nil || updatedCount.ModifiedCount == 0 {
				err = fmt.Errorf("UpdateAvatarInMembers %v:", err)
			}
			errChan <- err
		}()

		for i := 0; i < 7; i++ {
			if err := <-errChan; err != nil {
				log.Printf("Error updating avatar: %v\n", err)
			}
		}
		close(errChan)
	}()
}
