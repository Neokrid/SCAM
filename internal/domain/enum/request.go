package enum

type FriendRequestStatus string

const (
	Pending  FriendRequestStatus = "pending"
	Rejected FriendRequestStatus = "reject"
)

func (v FriendRequestStatus) String() string {
	return string(v)
}
