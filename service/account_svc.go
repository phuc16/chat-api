package service

import (
	"app/dto"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
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

func (s *AccountService) GetProfileByPhoneNumber(ctx context.Context, phoneNumber string) (res *entity.Account, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account, err := s.AccountRepo.SearchByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return
	}
	curAccount := entity.GetAccountFromContext(ctx)
	if curAccount.ID == account.ID {
		return account, nil
	}
	switch account.Setting.ShowBirthday {
	case entity.SHOW_BIRTHDAY_DM:
		account.Profile.Birthday = time.Date(0, account.Profile.Birthday.Month(), account.Profile.Birthday.Day(), account.Profile.Birthday.Hour(), account.Profile.Birthday.Minute(), account.Profile.Birthday.Second(), account.Profile.Birthday.Nanosecond(), account.Profile.Birthday.Location())
	case entity.SHOW_BIRTHDAY_NO:
		account.Profile.Birthday = time.Date(0, 0, 0, 0, 0, 0, 0, account.Profile.Birthday.Location())
	}
	return account, nil
}

func (s *AccountService) GetProfileByUserID(ctx context.Context, userID string) (res *entity.Account, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account, err := s.AccountRepo.SearchByUserID(ctx, userID)
	if err != nil {
		return
	}
	curAccount := entity.GetAccountFromContext(ctx)
	if curAccount.ID == account.ID {
		return account, nil
	}
	switch account.Setting.ShowBirthday {
	case entity.SHOW_BIRTHDAY_DM:
		account.Profile.Birthday = time.Date(0, account.Profile.Birthday.Month(), account.Profile.Birthday.Day(), account.Profile.Birthday.Hour(), account.Profile.Birthday.Minute(), account.Profile.Birthday.Second(), account.Profile.Birthday.Nanosecond(), account.Profile.Birthday.Location())
	case entity.SHOW_BIRTHDAY_NO:
		account.Profile.Birthday = time.Date(0, 0, 0, 0, 0, 0, 0, account.Profile.Birthday.Location())
	}
	return account, nil
}

func (s *AccountService) GetProfile(ctx context.Context) (res *entity.Account, err error) {
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
	s.UpdateAsyncSvc.UpdateAvatarAsync(ctx, account.Profile.Avatar, req.NewAvatar)
	return
}
