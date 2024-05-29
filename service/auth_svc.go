package service

import (
	"app/config"
	"app/dto"
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	AccountRepo IAccountRepo
	TokenRepo   ITokenRepo
	UserRepo    IUserRepo
}

func NewAuthService(accountRepo IAccountRepo, tokenRepo ITokenRepo, userRepo IUserRepo) *AuthService {
	return &AuthService{AccountRepo: accountRepo, TokenRepo: tokenRepo, UserRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, req *dto.AccountRegisterReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account := entity.Account{
		ID:          utils.NewID(),
		PhoneNumber: req.PhoneNumber,
		Pw:          utils.HashPassword(req.Password),
		Type:        entity.TYPE_ACCOUNT_PERSONAL,
		Profile: entity.Profile{
			UserID:     utils.NewID(),
			Avatar:     "https://grn-admin.mpoint.vn/uploads/avatar-mac-dinh.png",
			Background: "https://fptshop.com.vn/Uploads/Originals/2023/11/2/638345572887897897_anh-mac-dinh-zalo.jpg",
		},
		Role: entity.ROLE_USER,
		Setting: entity.Setting{
			AllowMessaging: entity.ALLOW_EVERYONE,
			ShowBirthday:   entity.SHOW_BIRTHDAY_DMY,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.AccountRepo.SaveAccount(ctx, &account)
	if err != nil {
		return
	}
	user := &entity.User{
		ID: account.Profile.UserID,
	}
	err = s.UserRepo.SaveUser(ctx, user)
	return
}

func (s *AuthService) Login(ctx context.Context, req *dto.AccountLoginReq) (accessToken string, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	account, err := s.AccountRepo.SearchByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return
	}
	if !utils.VerifyPassword(req.Password, account.Pw) {
		err = errors.PasswordIncorrect()
		return
	}
	accessToken, err = s.CreateToken(ctx, account)
	if err != nil {
		return
	}
	return
}

func (s *AuthService) CreateToken(ctx context.Context, account *entity.Account) (accessToken string, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	duration := time.Duration(int32(config.Cfg.HTTP.AccessTokenDuration)) * time.Minute
	appToken := &entity.Token{
		ID:          utils.NewID(),
		AccountID:   account.ID,
		PhoneNumber: account.PhoneNumber,
		UserID:      account.Profile.UserID,
		UserName:    account.Profile.UserName,
		Type:        entity.AccessTokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	err = s.TokenRepo.CreateToken(ctx, appToken)
	if err != nil {
		return
	}
	accessToken = appToken.SignedToken(config.Cfg.HTTP.Secret)
	return
}

func (s *AuthService) Logout(ctx context.Context, tokenStr string) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	token, err := entity.NewTokenFromEncoded(tokenStr, config.Cfg.HTTP.Secret)
	if err != nil {
		err = apperror.NewError(errors.CodeTokenError, err.Error())
		return
	}
	dbToken, err := s.TokenRepo.GetTokenByID(ctx, token.ID)
	if err != nil {
		return
	}
	err = s.TokenRepo.DeleteToken(ctx, dbToken)
	return
}

func (s *AuthService) Authenticate(ctx context.Context, tokenStr string) (res *entity.Account, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	token, err := entity.NewTokenFromEncoded(tokenStr, config.Cfg.HTTP.Secret)
	if err != nil {
		err = apperror.NewError(errors.CodeTokenError, err.Error())
		return
	}
	dbToken, err := s.TokenRepo.GetTokenByID(ctx, token.ID)
	if err != nil {
		return
	}
	dbAccount, err := s.AccountRepo.FindAccountByID(ctx, dbToken.AccountID)
	if err != nil {
		return
	}
	return dbAccount, nil
}
