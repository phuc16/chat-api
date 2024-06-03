package service

import (
	"context"
	"log"
	"sync"
)

type UpdateAsyncService struct {
	UserRepo  IUserRepo
	ChatRepo  IChatRepo
	GroupRepo IGroupRepo
}

func NewUpdateAsyncService(UserRepo IUserRepo, ChatRepo IChatRepo, GroupRepo IGroupRepo) *UpdateAsyncService {
	return &UpdateAsyncService{UserRepo: UserRepo, ChatRepo: ChatRepo, GroupRepo: GroupRepo}
}

func (s *UpdateAsyncService) UpdateAvatarAsync(ctx context.Context, userID, newAvatar string) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.UserRepo.UpdateAvatarInFriendRequest(ctx, userID, newAvatar)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.UserRepo.UpdateAvatarInConversation(ctx, userID, newAvatar)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.ChatRepo.UpdateAvatarInDelivery(ctx, userID, newAvatar)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.ChatRepo.UpdateAvatarInRead(ctx, userID, newAvatar)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.GroupRepo.UpdateAvatarInOwner(ctx, userID, newAvatar)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.GroupRepo.UpdateAvatarInAdmins(ctx, userID, newAvatar)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.GroupRepo.UpdateAvatarInMembers(ctx, userID, newAvatar)
		log.Println(err)
	}()

	wg.Wait()
}

func (s *UpdateAsyncService) UpdateNameAsync(ctx context.Context, userID, newName string) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.UserRepo.UpdateNameInFriendRequest(ctx, userID, newName)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.UserRepo.UpdateNameInConversation(ctx, userID, newName)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.ChatRepo.UpdateNameInDelivery(ctx, userID, newName)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.ChatRepo.UpdateNameInRead(ctx, userID, newName)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.GroupRepo.UpdateNameInOwner(ctx, userID, newName)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.GroupRepo.UpdateNameInAdmins(ctx, userID, newName)
		log.Println(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.GroupRepo.UpdateNameInMembers(ctx, userID, newName)
		log.Println(err)
	}()

	wg.Wait()
}
