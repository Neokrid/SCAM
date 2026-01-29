package auth

import (
	"context"
	"scam/internal/domain/dto/auth"
	"scam/internal/domain/dto/request"
	dto "scam/internal/domain/dto/response"
	"scam/internal/domain/dto/user"
	"scam/internal/domain/enum"
	"scam/internal/domain/services/token"
	apperrors "scam/internal/errors"
	userRepository "scam/internal/infrastructure/repository/user"
	"scam/pkg/applogger"
	"scam/pkg/trx"
	"time"

	"github.com/google/uuid"
)

var codeDelay time.Duration = time.Duration(time.Minute * 1)

type userService interface {
	CreateUserFromAuthCredentials(ctx context.Context, credintials request.RegisterCredentials) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string, password string) (*user.User, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, filter *userRepository.UserUpdateParams) error
}

type tokenService interface {
	GenerateUserTokens(ctx context.Context, id uuid.UUID) (*token.UserTokens, error)
	ParseToken(token string) (*token.CustomClaims, error)
	RefreshTokens(ctx context.Context, access, refresh string) (*token.UserTokens, error)
}

type smtpService interface {
	SendConfirmEmailCode(ctx context.Context, email string, action enum.EmailCodeAction) error
	ConfirmCode(ctx context.Context, email string, code string) (*auth.ConfirmationCode, error)
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	userService  userService
	smtpService  smtpService
	tokenService tokenService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	userService userService,
	smtpService smtpService,
	tokenService tokenService,
) *Service {
	return &Service{
		tx:           tx,
		logger:       logger,
		userService:  userService,
		smtpService:  smtpService,
		tokenService: tokenService,
	}
}

func (srv *Service) RegisterUser(ctx context.Context, credentials request.RegisterCredentials) (*dto.RegisterResponse, error) {
	var err error
	var user *user.User
	if err = srv.tx.Transaction(ctx, func(ctx context.Context) error {
		user, err = srv.userService.CreateUserFromAuthCredentials(ctx, credentials)
		if err != nil {
			return err
		}

		_, err = srv.SendConfirmationCode(ctx, request.LoginRequest{
			Email:    credentials.Email,
			Password: credentials.Password,
		}, enum.ConfirmCode)
		if err != nil {
			_ = ctx.Err().Error()
			return err
		}
		return err
	}); err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		UserId: user.Id,
	}, nil

}

func (srv *Service) SendConfirmationCode(ctx context.Context, req request.LoginRequest, action enum.EmailCodeAction) (*dto.SendCodeResponse, error) {
	_, err := srv.userService.GetUserByEmail(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &dto.SendCodeResponse{NextCodeDelay: codeDelay},
		srv.smtpService.SendConfirmEmailCode(ctx, req.Email, action)
}

func (srv *Service) ConfirmCode(ctx context.Context, req request.ConfimationCodeRequest) error {
	u, err := srv.userService.GetUserByEmail(ctx, req.Email, "")
	if err != nil {
		return err
	}
	code, err := srv.smtpService.ConfirmCode(ctx, req.Email, req.Code)
	if err != nil {
		return err
	}
	t := true
	switch code.Action {
	case enum.ConfirmCode:
		return srv.userService.UpdateUser(ctx, u.Id, &userRepository.UserUpdateParams{
			ConfirmedEmail: &t,
		})
	case enum.ForgotPassword:
		if req.NewPassword == "" {
			return apperrors.NoNewPassword
		}
		return srv.userService.UpdateUser(ctx, u.Id, &userRepository.UserUpdateParams{
			Password: &req.NewPassword,
		})
	default:
		return nil
	}
}

func (srv *Service) Login(ctx context.Context, req request.LoginRequest) (*token.UserTokens, error) {
	u, err := srv.userService.GetUserByEmail(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return srv.tokenService.GenerateUserTokens(ctx, u.Id)
}

func (srv *Service) RefreshTokens(ctx context.Context, req token.UserTokens) (*token.UserTokens, error) {
	return srv.tokenService.RefreshTokens(ctx, req.Access, req.Refresh)
}
