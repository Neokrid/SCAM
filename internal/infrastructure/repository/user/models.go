package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"fullName"`
	Status         string    `json:"status"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	ImgUrl         string    `json:"imgUrl"`
	ConfirmedEmail bool      `json:"confirmedEmail"`
	LastSeenAt     time.Time `json:"lastSeenAt"`
	BirthDate      time.Time `json:"birthDate"`
	CreatedAt      time.Time `json:"createdAt"`
}
