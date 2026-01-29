package social

import (
	"time"

	"github.com/google/uuid"
)

type FriendRequest struct {
	Id          uuid.UUID
	Sender_id   uuid.UUID
	Receiver_id uuid.UUID
	Status      string
	CreatedAt   time.Time
}

type Friendship struct {
	Id        uuid.UUID
	User_id   uuid.UUID
	Friend_id uuid.UUID
	Status    bool
	CreatedAt time.Time
}
