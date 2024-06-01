package service

import (
	"app/dto"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"math"
	"time"
)

type AccountService struct {
	UserRepo       IUserRepo
	AccountRepo    IAccountRepo
	UpdateAsyncSvc IUpdateAsyncSvc
}

func NewAccountService(userRepo IUserRepo, accountRepo IAccountRepo, updateAsyncSvc IUpdateAsyncSvc) *AccountService {
	return &AccountService{UserRepo: userRepo, AccountRepo: accountRepo, UpdateAsyncSvc: updateAsyncSvc}
}

func (s *AccountService) GetProfile(ctx context.Context) (res *entity.Profile, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account, err := s.AccountRepo.SearchByPhoneNumber(ctx, entity.GetAccountFromContext(ctx).PhoneNumber)
	if err != nil {
		return
	}
	return &account.Profile, nil
}

func (s *AccountService) GetProfileByPhoneNumber(ctx context.Context, phoneNumber string) (res *entity.Profile, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account, err := s.AccountRepo.SearchByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return
	}
	curAccount := entity.GetAccountFromContext(ctx)
	if curAccount.ID == account.ID {
		return &account.Profile, nil
	}
	switch account.Setting.ShowBirthday {
	case entity.SHOW_BIRTHDAY_DM:
		account.Profile.Birthday = time.Date(0, account.Profile.Birthday.Month(), account.Profile.Birthday.Day(), account.Profile.Birthday.Hour(), account.Profile.Birthday.Minute(), account.Profile.Birthday.Second(), account.Profile.Birthday.Nanosecond(), account.Profile.Birthday.Location())
	case entity.SHOW_BIRTHDAY_NO:
		account.Profile.Birthday = time.Date(0, 0, 0, 0, 0, 0, 0, account.Profile.Birthday.Location())
	}
	err = s.UpdateRecentSearchProfiles(ctx, curAccount.Profile.UserID, account.Profile)
	if curAccount.ID == account.ID {
		return nil, err
	}
	return &account.Profile, nil
}

func (s *AccountService) GetProfileByUserID(ctx context.Context, userID string) (res *entity.Profile, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account, err := s.AccountRepo.SearchByUserID(ctx, userID)
	if err != nil {
		return
	}
	curAccount := entity.GetAccountFromContext(ctx)
	if curAccount.ID == account.ID {
		return &account.Profile, nil
	}
	switch account.Setting.ShowBirthday {
	case entity.SHOW_BIRTHDAY_DM:
		account.Profile.Birthday = time.Date(0, account.Profile.Birthday.Month(), account.Profile.Birthday.Day(), account.Profile.Birthday.Hour(), account.Profile.Birthday.Minute(), account.Profile.Birthday.Second(), account.Profile.Birthday.Nanosecond(), account.Profile.Birthday.Location())
	case entity.SHOW_BIRTHDAY_NO:
		account.Profile.Birthday = time.Date(0, 0, 0, 0, 0, 0, 0, account.Profile.Birthday.Location())
	}
	return &account.Profile, nil
}

func (s *AccountService) GetAccountProfile(ctx context.Context) (res *entity.Account, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.AccountRepo.FindAccountByID(ctx, entity.GetAccountFromContext(ctx).ID)
}

func (s *AccountService) CheckPhoneNumber(ctx context.Context, req *dto.AccountCheckPhoneNumberReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.AccountRepo.SearchByPhoneNumber(ctx, req.PhoneNumber)
}

func (s *AccountService) ResetPassword(ctx context.Context, req *dto.AccountResetPasswordReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	_, err = s.AccountRepo.SearchByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return
	}
	_, err = s.AccountRepo.ChangePassword(ctx, req.PhoneNumber, req.NewPassword)
	if err != nil {
		return
	}
	return
}

func (s *AccountService) ChangePassword(ctx context.Context, req *dto.AccountChangePasswordReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account, err := s.AccountRepo.FindAccountByID(ctx, entity.GetAccountFromContext(ctx).ID)
	if err != nil {
		return
	}
	if !utils.VerifyPassword(req.CurPassword, account.Pw) {
		err = errors.PasswordIncorrect()
		return
	}
	_, err = s.AccountRepo.ChangePassword(ctx, account.PhoneNumber, req.NewPassword)
	if err != nil {
		return
	}
	return
}

func (s *AccountService) ChangeAvatar(ctx context.Context, req *dto.AccountChangeAvatarReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account, err := s.AccountRepo.FindAccountByID(ctx, entity.GetAccountFromContext(ctx).ID)
	if err != nil {
		return
	}
	newProfile := account.Profile
	newProfile.Avatar = req.NewAvatar
	_, err = s.AccountRepo.ChangeAvatar(ctx, account.PhoneNumber, newProfile)
	if err != nil {
		return
	}
	s.UpdateAsyncSvc.UpdateAvatarAsync(ctx, account.Profile.UserID, req.NewAvatar)
	return
}

func (s *AccountService) ChangeProfile(ctx context.Context, req *dto.AccountChangeProfileReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account, err := s.AccountRepo.FindAccountByID(ctx, entity.GetAccountFromContext(ctx).ID)
	if err != nil {
		return
	}
	account.Profile.UserName = req.UserName
	account.Profile.Gender = req.Gender
	account.Profile.Birthday = req.Birthday
	err = s.AccountRepo.UpdateAccount(ctx, account)
	if err != nil {
		return
	}
	s.UpdateAsyncSvc.UpdateNameAsync(ctx, account.Profile.UserID, req.UserName)
	return
}

func (s *AccountService) UpdateRecentSearchProfiles(ctx context.Context, userID string, newProfile entity.Profile) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	recentSearchProfiles, err := s.UserRepo.GetAllRecentSearchProfiles(ctx, userID)
	if err != nil {
		return err
	}
	var newRecentSearchProfiles []entity.Profile
	for _, profile := range recentSearchProfiles {
		if profile.UserID != newProfile.UserID {
			newRecentSearchProfiles = append(newRecentSearchProfiles, profile)
		}
	}
	newRecentSearchProfiles = append(newRecentSearchProfiles, newProfile)
	err = s.UserRepo.UpdateRecentSearchProfiles(ctx, userID, newRecentSearchProfiles)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountService) GetRecentSearchProfiles(ctx context.Context) (res []entity.Profile, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	curAccount := entity.GetAccountFromContext(ctx)
	recentSearchProfiles, err := s.UserRepo.GetAllRecentSearchProfiles(ctx, curAccount.Profile.UserID)

	startIndex := 0
	if len(recentSearchProfiles) > 3 {
		startIndex = len(recentSearchProfiles) - 3
	}
	return s.ReverseProfileSlice(recentSearchProfiles[startIndex:]), nil
}

func (s *AccountService) ReverseProfileSlice(slice []entity.Profile) []entity.Profile {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func (s *AccountService) GetSuggestFriendProfiles(ctx context.Context) (res []*entity.Profile, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	curAccount := entity.GetAccountFromContext(ctx)
	accounts, err := s.AccountRepo.GetAllAccounts(ctx)
	if err != nil {
		return
	}
	for i, account := range accounts {
		if curAccount.ID == account.ID {
			curAccount = *account
		} else {
			switch account.Setting.ShowBirthday {
			case entity.SHOW_BIRTHDAY_DM:
				account.Profile.Birthday = time.Date(0, account.Profile.Birthday.Month(), account.Profile.Birthday.Day(), account.Profile.Birthday.Hour(), account.Profile.Birthday.Minute(), account.Profile.Birthday.Second(), account.Profile.Birthday.Nanosecond(), account.Profile.Birthday.Location())
			case entity.SHOW_BIRTHDAY_NO:
				account.Profile.Birthday = time.Date(0, 0, 0, 0, 0, 0, 0, account.Profile.Birthday.Location())
			}
		}
		accounts[i] = account
	}
	// Find nearest neighbors of the target user
	var curAccountFriends []string
	for _, conversation := range curAccount.User.Conversations {
		if conversation.Type == entity.TYPE_FRIEND {
			curAccountFriends = append(curAccountFriends, conversation.IDUserOrGroup)
		}
	}
	nearestNeighbors := s.FindNearestNeighbors(&curAccount, accounts)
	res = []*entity.Profile{}
	// Suggest friends of nearest neighbors
	for _, neighbor := range nearestNeighbors {
		var neighborFriends []entity.Conversation
		for _, conversation := range neighbor.User.Conversations {
			if conversation.Type == entity.TYPE_FRIEND {
				neighborFriends = append(neighborFriends, conversation)
			}
		}
		for _, friend := range neighborFriends {
			if !s.Contains(curAccountFriends, friend.IDUserOrGroup) && friend.IDUserOrGroup != curAccount.User.ID {
				isSentRequest := false
				for _, conversation := range curAccount.User.Conversations {
					if conversation.IDUserOrGroup == friend.IDUserOrGroup && conversation.Type == entity.TYPE_REQUESTS {
						isSentRequest = true
						break
					}
				}
				if isSentRequest {
					continue
				}
				res = append(res, &entity.Profile{
					UserID:   friend.IDUserOrGroup,
					UserName: friend.ChatName,
					Avatar:   friend.ChatAvatar,
				})
			}
		}
	}
	return res, nil
}

// Function to find nearest neighbors based on cosine similarity
func (s *AccountService) FindNearestNeighbors(targetAccount *entity.Account, allAccounts []*entity.Account) []*entity.Account {
	// Calculate similarity with all users
	var similarities []float64
	for _, account := range allAccounts {
		similarity := s.CosineSimilarity(targetAccount, account)
		similarities = append(similarities, similarity)
	}

	// Find nearest neighbors
	var nearestNeighbors []*entity.Account
	for i, similarity := range similarities {
		if similarity > 0.8 && allAccounts[i].ID != targetAccount.ID { // Adjust threshold as needed
			nearestNeighbors = append(nearestNeighbors, allAccounts[i])
		}
	}

	return nearestNeighbors
}

// Function to calculate cosine similarity between two vectors
func (s *AccountService) CosineSimilarity(targetAccount, account *entity.Account) float64 {
	a, b := []string{}, []string{}
	for _, conversation := range targetAccount.User.Conversations {
		if conversation.Type == entity.TYPE_FRIEND {
			a = append(a, conversation.IDUserOrGroup)
		}
	}
	for _, conversation := range account.User.Conversations {
		if conversation.Type == entity.TYPE_FRIEND {
			b = append(b, conversation.IDUserOrGroup)
		}
	}
	// Convert friend lists to binary vectors
	length := len(a)
	if len(b) > length {
		length = len(b)
	}
	vectorA := make([]int, length)
	vectorB := make([]int, length)
	for i, friendID := range a {
		if s.Contains(b, friendID) {
			vectorA[i] = 1
		}
	}
	for i, friendID := range b {
		if s.Contains(a, friendID) {
			vectorB[i] = 1
		}
	}

	// Calculate dot product
	dotProduct := 0.0
	for i := 0; i < len(vectorA); i++ {
		dotProduct += float64(vectorA[i] * vectorB[i])
	}

	// Calculate magnitudes
	magnitudeA := 0.0
	magnitudeB := 0.0
	for _, val := range vectorA {
		magnitudeA += float64(val * val)
	}
	for _, val := range vectorB {
		magnitudeB += float64(val * val)
	}
	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)

	// Calculate cosine similarity
	similarity := dotProduct / (magnitudeA * magnitudeB)
	return similarity
}

// Function to check if a user ID exists in a list of friend IDs
func (s *AccountService) Contains(friends []string, userID string) bool {
	for _, friendID := range friends {
		if friendID == userID {
			return true
		}
	}
	return false
}
