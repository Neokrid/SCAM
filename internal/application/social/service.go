package social

import (
	"context"
	"scam/internal/domain/dto/request"
	"scam/pkg/applogger"
	"scam/pkg/trx"

	"github.com/google/uuid"
)

type socialService interface {
	SendFriendRequest(ctx context.Context, targetId request.SendFriendRequest, senderID uuid.UUID) error
	RespondFriendRequest(ctx context.Context, requestID uuid.UUID, userID uuid.UUID, status string) error
	CancelFriendRequest(ctx context.Context, requestID, userID uuid.UUID) error
	GetRequests(ctx context.Context, userID uuid.UUID) (*request.RequestsResponse, error)
	GetFriendsList(ctx context.Context, userID uuid.UUID) ([]request.UserListRequest, error)
	RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error
	RemoveFollower(ctx context.Context, userID, followerID uuid.UUID) error
	GetFollowersList(ctx context.Context, userID uuid.UUID) ([]request.UserListRequest, error)
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	socialService socialService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	socialService socialService,
) *Service {
	return &Service{
		tx:            tx,
		logger:        logger,
		socialService: socialService,
	}
}

func (srv *Service) SendFriendRequest(ctx context.Context, targetId request.SendFriendRequest, senderID uuid.UUID) error {
	return srv.socialService.SendFriendRequest(ctx, targetId, senderID)
}

func (srv *Service) RespondFriendRequest(ctx context.Context, requestID uuid.UUID, userID uuid.UUID, status string) error {
	return srv.socialService.RespondFriendRequest(ctx, requestID, userID, status)
}

func (srv *Service) CancelFriendRequest(ctx context.Context, requestID, userID uuid.UUID) error {
	return srv.socialService.CancelFriendRequest(ctx, requestID, userID)
}

func (srv *Service) GetRequests(ctx context.Context, userID uuid.UUID) (*request.RequestsResponse, error) {
	req, err := srv.socialService.GetRequests(ctx, userID)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (srv *Service) GetFriendsList(ctx context.Context, userID uuid.UUID) ([]request.UserListRequest, error) {
	userList, err := srv.socialService.GetFriendsList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

func (srv *Service) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	return srv.socialService.RemoveFriend(ctx, userID, friendID)
}

func (srv *Service) RemoveFollower(ctx context.Context, userID, followerID uuid.UUID) error {
	return srv.socialService.RemoveFollower(ctx, userID, followerID)
}

func (srv *Service) GetFollowersList(ctx context.Context, userID uuid.UUID) ([]request.UserListRequest, error) {
	followersList, err := srv.socialService.GetFollowersList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return followersList, nil
}
