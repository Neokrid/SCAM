package social

import (
	"context"
	"scam/internal/domain/dto/user"
	apperrors "scam/internal/errors"
	"scam/internal/infrastructure/repository/common"
	"scam/pkg/database/postgres"

	"github.com/google/uuid"
)

type Repository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) *Repository {
	return &Repository{conn: conn}
}

func (repo *Repository) SendFriendRequest(ctx context.Context, friendRequest *FriendRequest) error {
	query := `
        INSERT INTO friend_requests (
            id, 
            sender_id, 
            receiver_id, 
            status, 
            created_at
        )
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := repo.conn.Exec(ctx, query, friendRequest.Id, friendRequest.Sender_id, friendRequest.Receiver_id, friendRequest.Status, friendRequest.CreatedAt)
	if err != nil {
		if common.IsUniqueErr(err) {
			return apperrors.NotUnique
		}
	}
	return err

}

func (repo *Repository) RespondFriendRequest(ctx context.Context, requestID uuid.UUID, userID uuid.UUID, status string) error {
	query := `
        UPDATE friend_requests 
        SET status = $1
        WHERE id = $2
    `
	res, err := repo.conn.Exec(ctx, query, status, requestID)
	if err != nil {
		return err
	}
	rows := res.RowsAffected()
	if rows == 0 {
		return apperrors.FriendRequestNotFound
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*FriendRequest, error) {
	query := `SELECT id, sender_id, receiver_id, status FROM friend_requests WHERE id = $1`
	var req FriendRequest
	err := r.conn.QueryRow(ctx, query, id).Scan(&req.Id, &req.Sender_id, &req.Receiver_id, &req.Status)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *Repository) CreateFriendship(ctx context.Context, friendShip *Friendship) error {

	upsertQuery := `
        INSERT INTO friendships (user_id, friend_id, status, created_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, friend_id) 
        DO UPDATE SET status = EXCLUDED.status
    `

	if friendShip.Status {
		if _, err := r.conn.Exec(ctx, upsertQuery, friendShip.Friend_id, friendShip.User_id, friendShip.Status, friendShip.CreatedAt); err != nil {
			return err
		}
		if _, err := r.conn.Exec(ctx, upsertQuery, friendShip.User_id, friendShip.Friend_id, friendShip.Status, friendShip.CreatedAt); err != nil {
			return err
		}
	} else {

		if _, err := r.conn.Exec(ctx, upsertQuery, friendShip.Friend_id, friendShip.User_id, friendShip.Status, friendShip.CreatedAt); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) CancelFriendRequest(ctx context.Context, requestID, userID uuid.UUID) error {
	query := `
        DELETE FROM friend_requests 
        WHERE id = $1 AND sender_id = $2 AND status = 'PENDING'
    `
	res, err := r.conn.Exec(ctx, query, requestID, userID)
	if err != nil {
		return err
	}
	rows := res.RowsAffected()
	if rows == 0 {
		return apperrors.FriendRequestNotFound
	}
	return nil
}

func (r *Repository) GetRequests(ctx context.Context, userID uuid.UUID) ([]user.FriendRequestWithUser, []user.FriendRequestWithUser, error) {
	incomingQuery := `
        SELECT fr.id, u.id, u.username, u.img_url, fr.status, fr.created_at
        FROM friend_requests fr
        JOIN users u ON fr.sender_id = u.id
        WHERE fr.receiver_id = $1 AND fr.status = 'pending'
    `
	outgoingQuery := `
        SELECT fr.id, u.id, u.username, u.img_url, fr.status, fr.created_at
        FROM friend_requests fr
        JOIN users u ON fr.receiver_id = u.id
        WHERE fr.sender_id = $1 AND fr.status = 'pending'
    `
	incoming, err := r.fetchRequests(ctx, incomingQuery, userID)
	if err != nil {
		return nil, nil, err
	}
	outgoing, err := r.fetchRequests(ctx, outgoingQuery, userID)
	if err != nil {
		return nil, nil, err
	}
	return incoming, outgoing, nil
}

func (r *Repository) fetchRequests(ctx context.Context, query string, arg uuid.UUID) ([]user.FriendRequestWithUser, error) {
	rows, err := r.conn.Query(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []user.FriendRequestWithUser
	for rows.Next() {
		var i user.FriendRequestWithUser
		if err := rows.Scan(&i.ID, &i.UserID, &i.Username, &i.ImgURL, &i.Status, &i.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, i)
	}
	return list, nil
}

func (r *Repository) GetFriendsList(ctx context.Context, userID uuid.UUID) ([]user.User, error) {
	query := `
        SELECT u.id, u.username, u.email, u.img_url, u.created_at
        FROM friendships f
        JOIN users u ON f.friend_id = u.id
        WHERE f.user_id = $1 AND f.status = true
    `
	rows, err := r.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	var friendsList []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.ImgUrl, &u.CreatedAt); err != nil {
			return nil, err
		}

		friendsList = append(friendsList, u)
	}

	return friendsList, nil
}

func (r *Repository) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	_, err := r.conn.Exec(ctx, `DELETE FROM friendships WHERE user_id = $1 AND friend_id = $2`, userID, friendID)
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(ctx, `UPDATE friendships SET status = false WHERE user_id = $1 AND friend_id = $2`, friendID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) RemoveFollower(ctx context.Context, userID, followerID uuid.UUID) error {
	query := `DELETE FROM friendships WHERE user_id = $1 AND friend_id = $2 AND status = false`

	_, err := r.conn.Exec(ctx, query, followerID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetFollowersList(ctx context.Context, userID uuid.UUID) ([]user.User, error) {
	query := `
        SELECT u.id, u.username, u.email, u.img_url, u.created_at
        FROM friendships f
        JOIN users u ON f.friend_id = u.id
        WHERE f.user_id = $1 AND f.status = false
    `
	rows, err := r.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	var friendsList []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.ImgUrl, &u.CreatedAt); err != nil {
			return nil, err
		}

		friendsList = append(friendsList, u)
	}

	return friendsList, nil
}
