package social

import (
	"context"
	"scam/internal/domain/dto/request"
	"scam/pkg/apperror"
	"scam/pkg/applogger"
	"scam/pkg/constants"
	"scam/pkg/response"
	"scam/pkg/util"

	"github.com/gin-gonic/gin"
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

type Controller struct {
	lgr     applogger.Logger
	builder *response.Builder

	socialService socialService
}

func NewController(
	logger applogger.Logger,
	builder *response.Builder,
	socialService socialService,
) *Controller {
	return &Controller{
		lgr:           logger,
		builder:       builder,
		socialService: socialService,
	}
}

func (h *Controller) Init(authApi *gin.RouterGroup) {
	friends := authApi.Group("/friends")
	{
		friends.POST("/request", h.sendFriendRequest)
		friends.POST("/request/:id/accept", h.acceptFriendRequest)
		friends.POST("/request/:id/reject", h.rejectFriendRequest)
		friends.DELETE("/request/:id/cancel", h.cancelFriendRequest)
		friends.GET("/requests", h.getFriendRequests)
		friends.GET("/list", h.getFriendsList)
		friends.DELETE("/list/:user_id", h.removeFriend)
		friends.GET("/followers", h.getFollowersList)
		friends.DELETE("/followers/:user_id", h.removeFollower)
	}
}

func (h *Controller) sendFriendRequest(c *gin.Context) {
	var targetId request.SendFriendRequest
	senderID, err := util.GetUserId(c.Request.Context())
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	err = c.BindJSON(&targetId)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	err = h.socialService.SendFriendRequest(c, targetId, senderID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(c, nil))
}

func (h *Controller) acceptFriendRequest(c *gin.Context) {
	ctx := c.Request.Context()
	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	userID, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	err = h.socialService.RespondFriendRequest(ctx, requestID, userID, "accept")
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(c, nil))
}

func (h *Controller) rejectFriendRequest(c *gin.Context) {
	ctx := c.Request.Context()
	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	userID, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	err = h.socialService.RespondFriendRequest(ctx, requestID, userID, "reject")
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(c, nil))
}

func (h *Controller) cancelFriendRequest(c *gin.Context) {
	ctx := c.Request.Context()
	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	userID, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	err = h.socialService.CancelFriendRequest(ctx, requestID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(c, nil))
}
func (h *Controller) getFriendRequests(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	resp, err := h.socialService.GetRequests(ctx, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, resp))
}
func (h *Controller) getFriendsList(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	friends, err := h.socialService.GetFriendsList(ctx, userID)
	if err != nil {
		_ = c.Error(apperror.NewInternalError(err))
		return
	}
	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, friends))
}

func (h *Controller) removeFriend(c *gin.Context) {
	friendIDStr := c.Param("user_id")
	friendID, err := uuid.Parse(friendIDStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	ctx := c.Request.Context()
	userID, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	err = h.socialService.RemoveFriend(ctx, userID, friendID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(c, nil))
}

func (h *Controller) getFollowersList(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	friends, err := h.socialService.GetFollowersList(ctx, userID)
	if err != nil {
		_ = c.Error(apperror.NewInternalError(err))
		return
	}
	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, friends))
}

func (h *Controller) removeFollower(c *gin.Context) {
	followerIdStr := c.Param("user_id")
	followerId, err := uuid.Parse(followerIdStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	ctx := c.Request.Context()
	userID, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	err = h.socialService.RemoveFollower(ctx, userID, followerId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, nil))
}
