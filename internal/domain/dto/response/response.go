package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterResponse struct {
	UserId uuid.UUID `json:"userId"`
}

type SendCodeResponse struct {
	NextCodeDelay time.Duration `json:"nextCodeDelay"`
}

type ChangePictureResponse struct {
	NewImgurl string `json:"newImgUrl"`
}

type ProfileResponse struct {
	UserId    uuid.UUID `json:"userId"`
	ImgUrl    string    `json:"imgUrl"`
	FullName  string    `json:"fullName"`
	Status    string    `json:"status"`
	Online    bool      `json:"online"`
	BirthDate time.Time `json:"birthDate"`
}
