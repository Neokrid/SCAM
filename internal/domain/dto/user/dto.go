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
