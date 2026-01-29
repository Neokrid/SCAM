package user

import (
	"context"
	"scam/internal/domain/dto/request"
	resp "scam/internal/domain/dto/response"
	"scam/internal/domain/dto/user"
	apperrors "scam/internal/errors"
	"scam/pkg/apperror"
	"scam/pkg/applogger"
	"scam/pkg/constants"
	"scam/pkg/response"
	"scam/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userService interface {
	ChangeProfilePicture(ctx context.Context, req request.ChangeProfilePicture, host string) (*resp.ChangePictureResponse, error)
	GetUserById(ctx context.Context, userId uuid.UUID, host string) (*user.User, error)
	UpdateMyProfile(ctx context.Context, userId uuid.UUID, req request.UpdateProfileRequest) error
}

type Controller struct {
	lgr     applogger.Logger
	builder *response.Builder

	userService userService
}

func NewController(logger applogger.Logger, builder *response.Builder, userService userService) *Controller {
	return &Controller{
		lgr:     logger,
		builder: builder,

		userService: userService,
	}
}

func (h *Controller) Init(api, authApi *gin.RouterGroup) {
	user := api.Group("/user")
	userAuth := authApi.Group("/user")
	{
		userAuth.PUT("/picture", h.changeProfilePicture)
		userAuth.GET("me", h.getMyProfile)
		userAuth.PUT("me", h.updateMyProfile)
		user.GET("/profile/:id", h.getUserById)
	}
}

// @Summary change_profile_picture
// @Description сменить аватарку пользователя
// @Tags user
// @Produce json
// @Param data body request.ChangeProfilePicture true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{data=resp.ChangePictureResponse}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: user_not_found"
// @Router /rl/api/v1/user/register [post]
func (h *Controller) changeProfilePicture(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.ChangeProfilePicture
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}
	req.UserId = userId
	picUrl, err := h.userService.ChangeProfilePicture(ctx, req, c.Request.Host)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, picUrl))
}

// @Summary get_user_by_id
// @Description получить юзера по айди
// @Tags user
// @Produce json
// @Param path id string true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Success 200 {object} response.Response{data=resp.ChangePictureResponse}
// @Failure 400 {object} response.Response{} "possible codes: bind_path, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: user_not_found"
// @Router /rl/api/v1/user/profile/{id} [post]
func (h *Controller) getUserById(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}

	user, err := h.userService.GetUserById(ctx, id, c.Request.Host)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, user))
}

// @Summary get_my_profile
// @Description get my profile
// @Tags user
// @Produce json
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{data=user.User}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Router /rl/api/v1/user/me [get]
func (h *Controller) getMyProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	user, err := h.userService.GetUserById(ctx, userId, c.Request.Host)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, user))
}

// @Summary update_my_profile
// @Description обновляет мой профиль
// @Tags user
// @Produce json
// @Param data body request.UpdateProfileRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{data=user.User}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, bind_body"
// @Router /rl/api/v1/user/me [put]
func (h *Controller) updateMyProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}
	var req request.UpdateProfileRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}

	err = h.userService.UpdateMyProfile(ctx, userId, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := h.userService.GetUserById(ctx, userId, c.Request.Host)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, user))
}
