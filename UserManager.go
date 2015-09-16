package main

import (
	"fmt"

	"gopkg.in/jmcvetta/napping.v1"
)

// UserResponse list of users
type UserResponse struct {
	Ok      bool
	Members []Member
}

// Member single member
type Member struct {
	ID       string
	Name     string
	Deleted  bool
	RealName string `json:"real_name"`
	Profile  Profile
}

// Profile user profile
type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string
	Image192  string `json:"image_192"`
}

var users []Member
var reqParams = &napping.Params{"token": config.SlackToken}

func loadUsers() {
	var err interface{}
	var userResponse UserResponse

	napping.Get("https://slack.com/api/users.list", reqParams, &userResponse, err)

	fmt.Printf("Loaded %d users\n", len(userResponse.Members))

	users = userResponse.Members
}
