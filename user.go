package main

import (
	"fmt"
	"log"

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

func loadUsers() {
	var err interface{}
	var userResponse UserResponse
	var reqParams = &napping.Params{"token": config.SlackToken}

	napping.Get("https://slack.com/api/users.list", reqParams, &userResponse, err)

	log.Printf("Loaded %d users\n", len(userResponse.Members))

	users = userResponse.Members
}

func findMemberByTag(tag string) (Member, error) {
	for _, user := range users {
		if user.Name == tag {
			return user, nil
		}
	}

	return Member{}, fmt.Errorf("Member with tag %s could not be found!", tag)
}

func findMemberByID(ID string) (Member, error) {
	for _, user := range users {
		if user.ID == ID {
			return user, nil
		}
	}

	return Member{}, fmt.Errorf("Member with ID %s could not be found!", ID)
}
