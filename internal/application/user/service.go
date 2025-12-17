package auth

import (
	"context"
	"mime/multipart"

	"scam/internal/domain/dto/request"
	respDto "scam/internal/domain/dto/response"
	userDto "scam/internal/domain/dto/user"
	"scam/internal/infrastructure/repository/user"
	"scam/pkg/applogger"
	"scam/pkg/trx"

	"github.com/google/uuid"
)

type userService interface {
	UpdateUser(ctx context.Context, userId uuid.UUID, filter *user.UserUpdateParams) error
	GetUserById(ctx context.Context, userId uuid.UUID, password string) (*userDto.User, error)
}

type fileService interface {
	NewFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	userService userService
	fileService fileService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	userService userService,
	fileService fileService,
) *Service {
	return &Service{
		tx:          tx,
		logger:      logger,
		userService: userService,
		fileService: fileService,
	}
}

// todo add reg exp check for password and username and email

func (srv *Service) ChangeProfilePicture(ctx context.Context, req request.ChangeProfilePicture, host string) (*respDto.ChangePictureResponse, error) {
	filename, err := srv.fileService.NewFile(ctx, req.File)
	if err != nil {
		return nil, err
	}
	err = srv.userService.UpdateUser(ctx, req.UserId, &user.UserUpdateParams{
		ImgUrl: &filename,
	})
	return &respDto.ChangePictureResponse{
		NewImgurl: host + "/statics/images/" + filename,
	}, err
}

func (srv *Service) GetUserById(ctx context.Context, userId uuid.UUID, host string) (*userDto.User, error) {
	u, err := srv.userService.GetUserById(ctx, userId, "")
	if err != nil {
		return nil, err
	}
	u.ImgUrl = host + "/statics/images/" + u.ImgUrl
	return u, nil
}
