package dbmodels

import "time"

// MessageTime is the data model for the mongodb saved messages
type MessageTime struct {
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Content  string    `json:"content"`
	Date     time.Time `json:"time"`
}

// CropSizes ..
type CropSizes struct {
	Y      int `json:"y"`
	X      int `json:"x"`
	Width  int `json:"width"`
	Height int `json:"height"`
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

// NameAndImage request model
type NameAndImage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// NameAndImageAndAreaToCrop ..
type NameAndImageAndAreaToCrop struct {
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	AreaToCrop CropSizes `json:"areatocrop"`
}

// User ...
type User struct {
	Name           string           `json:"name"`
	Pass           string           `json:"pass"`
	Friends        []string         `json:"friends"`
	FriendRequests NameAndDateArray `json:"friendrequests"`
	ProfileImage   string           `json:"profileimage"`
}

// Name ...
type Name struct {
	Name string `json:"name"`
}

// Image ...
type Image struct {
	ProfileImage string `json:"profileimage"`
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
