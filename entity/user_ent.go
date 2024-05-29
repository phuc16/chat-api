package entity

type User struct {
	ID             string          `bson:"id" json:"id"`
	FriendRequests []FriendRequest `bson:"friend_requests" json:"friendRequests"`
	Conversations  []Conversation  `bson:"conversations" json:"conversations"`
}
