package dbmodels

import "time"

// MessageModel is the data model for the mongodb saved messages
type MessageModel struct {
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Content  string    `json:"content"`
	Date     time.Time `json:"time"`
}

// MessageRequest is the data model for the request messages
type MessageRequest struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

// UserModel is the data model for the user
type UserModel struct {
	Name    string   `json:"name"`
	Friends []string `json:"friends"`
}

// FriendRequest request model
type FriendRequest struct {
	User      string `json:"user"`
	NewFriend string `json:"newFriend"`
}

// FriendResult ...
type FriendResult struct {
	Name           string            `json:"name"`
	Pass           string            `json:"pass"`
	Friends        []string          `json:"friends"`
	FriendRequests FriendAddRequests `json:"friendRequests"`
}

// GetFriendsRequest ...
type GetFriendsRequest struct {
	Name string `json:"name"`
}

// GetWithFilterRequest ...
type GetWithFilterRequest struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

// LoginRequest ...
type LoginRequest struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

// FriendAddRequests ..
type FriendAddRequests []FriendAddRequest

// FriendAddRequest ..
type FriendAddRequest struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type GetFriendsRequestsResult struct {
	FriendRequests []FriendAddRequest `json:"friendRequests"`
}
