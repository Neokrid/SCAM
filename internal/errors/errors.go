package apperrors

import (
	"scam/pkg/apperror"
)

var (
	InvalidAuthorizationHeader = apperror.NewUnauthorizedError("invalid authorization header", "invalid_authorization_header")
	InvalidTokenError          = apperror.NewUnauthorizedError("invalid token", "invalid_token")

	UserNotFound           = apperror.NewInvalidDataError("user not found", "user_not_found")
	IncorrectPassword      = apperror.NewUnauthorizedError("incorrect password", "incorrect_password")
	ConfirmCodeAlreadySend = apperror.NewInvalidDataError("confirm code already send", "confirm_code_already_send")
	ConfirmCodeNotExist    = apperror.NewInvalidDataError("confirm code not exist", "confirm_code_not_exist")
	ConfirmCodeIncorrect   = apperror.NewInvalidDataError("confirm code incorrect", "confirm_code_incorrect")

	TokenClaimsError = apperror.NewInvalidDataError("bad token claims", "bad_token_claims")
	TokensDontMatch  = apperror.NewInvalidDataError("tokens dont match", "tokens_dont_match")
	TokenDontExist   = apperror.NewInvalidDataError("token dont exist", "token_dont_exist")

	NoNewPassword = apperror.NewBadRequestError("no new password", "no_new_password")
	NotUnique     = apperror.NewInvalidDataError("not unique", "not_unique")

	NoPermissionsRequest  = apperror.NewBadRequestError("no permission to respond to this request", "no_permission_to_respond_to_this_request")
	InvalidStatus         = apperror.NewBadRequestError("invalid status", "invalid_status")
	FriendRequestNotFound = apperror.NewInvalidDataError("friend request not found", "friend_request_not_found")
)

// коды динамических ошибок:
// bad_refresh_token 422
// invalid_X-Request-Id 400
// bad_access_token 422
// bind_path 400
