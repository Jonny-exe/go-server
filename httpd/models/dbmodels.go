package dbmodels

import "time"

// MessageModel is the data model for the mongodb saved messages
type MessageModel struct {
	Sender   string    `json:"Sender"`
	Receiver string    `json:"Receiver"`
	Content  string    `json:"Content"`
	Date     time.Time `json:"Time"`
}

// MessageRequest is the data model for the request messages
type MessageRequest struct {
	Sender   string `json:"Sender"`
	Receiver string `json:"Receiver"`
	Content  string `json:"Content"`
}

// UserModel is the data model for the user
type UserModel struct {
	Name    string   `json:"Name"`
	Friends []string `json:"Friends"`
}

// FriendRequest request model
type FriendRequest struct {
	User      string `json:"User"`
	NewFriend string `json:"NewFriend"`
}

// FriendResult ...
type FriendResult struct {
	User    string   `json:"User"`
	Friends []string `json:"Friends"`
}
