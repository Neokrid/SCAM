package social

import (
	"context"
	"scam/internal/domain/dto/request"
	"scam/internal/domain/dto/user"
	apperrors "scam/internal/errors"
	"scam/internal/infrastructure/repository/social"
	"scam/pkg/applogger"
	"scam/pkg/constants"
	"scam/pkg/trx"
	"scam/pkg/util"

	"github.com/google/uuid"
)

type socialRepo interface {
	SendFriendRequest(ctx context.Context, friendRequest *social.FriendRequest) error
	RespondFriendRequest(ctx context.Context, requestID uuid.UUID, userID uuid.UUID, status string) error
	GetByID(ctx context.Context, id uuid.UUID) (*social.FriendRequest, error)
	CreateFriendship(ctx context.Context, friendShip *social.Friendship) error
	CancelFriendRequest(ctx context.Context, requestID, userID uuid.UUID) error
	GetRequests(ctx context.Context, userID uuid.UUID) ([]user.FriendRequestWithUser, []user.FriendRequestWithUser, error)
	GetFriendsList(ctx context.Context, userID uuid.UUID) ([]user.User, error)
	RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error
	RemoveFollower(ctx context.Context, userID, followerID uuid.UUID) error
	GetFollowersList(ctx context.Context, userID uuid.UUID) ([]user.User, error)
}

type Service struct {
	tx         trx.TransactionManager
	logger     applogger.Logger
	socialRepo socialRepo
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	socialRepo socialRepo,
) *Service {
	return &Service{
		tx:         tx,
		logger:     logger,
		socialRepo: socialRepo,
	}
}

func (srv *Service) SendFriendRequest(ctx context.Context, targetId request.SendFriendRequest, senderID uuid.UUID) error {
	friendRequest := social.FriendRequest{
		Id:          util.NewUUID(),
		Sender_id:   senderID,
		Receiver_id: targetId.TargetId,
		Status:      constants.FriendRequestStatus,
		CreatedAt:   util.GetCurrentUTCTime(),
	}

	return srv.socialRepo.SendFriendRequest(ctx, &friendRequest)
}

func (srv *Service) RespondFriendRequest(ctx context.Context, requestID uuid.UUID, userID uuid.UUID, status string) error {

	req, err := srv.socialRepo.GetByID(ctx, requestID)
	friendShip := social.Friendship{
		Id:        util.NewUUID(),
		User_id:   userID,
		Friend_id: req.Sender_id,
		CreatedAt: util.GetCurrentUTCTime(),
	}
	if err != nil {
		return err
	}
	if req.Receiver_id != userID {
		return apperrors.NoPermissionsRequest
	}
	var newStatus string
	switch status {
	case "accept":
		newStatus = "accept"
		friendShip.Status = true
		err = srv.socialRepo.CreateFriendship(ctx, &friendShip)
		if err != nil {
			return err
		}
	case "reject":
		newStatus = "reject"
		friendShip.Status = false
		err = srv.socialRepo.CreateFriendship(ctx, &friendShip)
		if err != nil {
			return err
		}
	default:
		return apperrors.InvalidStatus
	}

	return srv.socialRepo.RespondFriendRequest(ctx, requestID, userID, newStatus)
}

func (srv *Service) CancelFriendRequest(ctx context.Context, requestID, userID uuid.UUID) error {
	return srv.socialRepo.CancelFriendRequest(ctx, requestID, userID)
}

func (srv *Service) GetRequests(ctx context.Context, userID uuid.UUID) (*request.RequestsResponse, error) {
	in, out, err := srv.socialRepo.GetRequests(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &request.RequestsResponse{
		Incoming: in,
		Outgoing: out,
	}, nil
}

func (srv *Service) GetFriendsList(ctx context.Context, userID uuid.UUID) ([]request.UserListRequest, error) {
	userList, err := srv.socialRepo.GetFriendsList(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []request.UserListRequest
	for _, u := range userList {
		result = append(result, request.UserListRequest{
			Id:       u.Id,
			Username: u.Username,
			ImgUrl:   u.ImgUrl,
		})
	}
	return result, nil
}

func (srv *Service) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	return srv.socialRepo.RemoveFriend(ctx, userID, friendID)
}

func (srv *Service) RemoveFollower(ctx context.Context, userID, followerID uuid.UUID) error {
	return srv.socialRepo.RemoveFollower(ctx, userID, followerID)
}

func (srv *Service) GetFollowersList(ctx context.Context, userID uuid.UUID) ([]request.UserListRequest, error) {
	followersList, err := srv.socialRepo.GetFriendsList(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []request.UserListRequest
	for _, u := range followersList {
		result = append(result, request.UserListRequest{
			Id:       u.Id,
			Username: u.Username,
			ImgUrl:   u.ImgUrl,
		})
	}
	return result, nil
}
