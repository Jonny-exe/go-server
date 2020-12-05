package dbmodels

import "time"

// MessageTime is the data model for the mongodb saved messages
type MessageTime struct {
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Content  string    `json:"content"`
	Date     time.Time `json:"time"`
}

// Message is the data model for the request messages
type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

// NameAndNewFriend request model
type NameAndNewFriend struct {
	Name      string `json:"name"`
	NewFriend string `json:"newFriend"`
}

// User ...
type User struct {
	Name           string           `json:"name"`
	Pass           string           `json:"pass"`
	Friends        []string         `json:"friends"`
	FriendRequests NameAndDateArray `json:"friendrequests"`
}

// Name ...
type Name struct {
	Name string `json:"name"`
}

// SenderAndReceiver ...
type SenderAndReceiver struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

// NameAndPass ...
type NameAndPass struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

// NameAndDateArray ..
type NameAndDateArray []NameAndDate

// NameAndDate ..
type NameAndDate struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

// NameAndDateStruct ...
type NameAndDateStruct struct {
	FriendRequests []NameAndDate `json:"friendrequests"`
}

// Friends ..
type Friends struct {
	Friends []string `json:"friends"`
}

// NameAndFriendToRemove ...
type NameAndFriendToRemove struct {
	Name           string `json:"name"`
	FriendToRemove string `json:"friendtoremove"`
}
