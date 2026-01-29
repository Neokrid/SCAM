package user

import (
	"scam/internal/infrastructure/repository/user"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	ImgUrl    string    `json:"imgUrl"`
	Status    string    `json:"status"`
	BirthDate time.Time `json:"birthDate"`
	CreatedAt time.Time `json:"createdAt"`
}

func UserDtoFromEntity(entity *user.User) *User {
	return &User{
		Id:        entity.Id,
		Username:  entity.Username,
		FullName:  entity.FullName,
		Email:     entity.Email,
		ImgUrl:    entity.ImgUrl,
		Status:    entity.Status,
		BirthDate: entity.BirthDate,
		CreatedAt: entity.CreatedAt,
	}
}

type FriendRequestWithUser struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	ImgURL    string    `json:"img_url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
