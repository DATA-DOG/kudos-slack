package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Kudo main kudos
type Kudo struct {
	ID         int64
	Kudo       string
	MemberTo   *Member
	MemberFrom *Member
	LikeCount  int
}

var kudos []Kudo

func main() {
	readConfig()
	loadUsers()
	loadDatabase()

	router := httprouter.New()
	router.GET("/", index)

	router.POST("/kudo", handleKudoCmd)

	fmt.Print("Listening on port ", config.Port, "...")
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", config.Port), router))
}

func handleKudoCmd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	accessToken := r.PostFormValue("token")

	if accessToken != config.SlackCommandToken {
		fmt.Fprint(w, "Invalid access token")
		return
	}

	memberFrom, err := getMemberFrom(r)
	if err != nil {
		fmt.Fprint(w, "Invalid user provided")
		return
	}

	command, target, extra, err := getCommandParams(r)
	if err != nil {
		printKudoUsage(w)
		return
	}

	switch command {
	case "to":
		handleNewKudoCommand(w, memberFrom, target, extra)
	case "like":
		handleLikeCommand(w, memberFrom, target)
	default:
		printKudoUsage(w)
	}
}

func printKudoUsage(w http.ResponseWriter) {
	fmt.Fprint(w, "New kudo: `/kudos to @user reason`\nLike user latest kudo: `/kudos like @user`")
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	for _, kudo := range kudos {
		fmt.Fprint(w, kudo.MemberFrom.Name, ": ", kudo.MemberTo.Name, ", ", kudo.Kudo, kudo.LikeCount, "\n")
	}
}

func getMemberFrom(r *http.Request) (*Member, error) {
	memberFromTag := r.PostFormValue("user_name")

	var memberFrom, err = findMemberByTag(memberFromTag)
	if err != nil {
		loadUsers()
		memberFrom, err = findMemberByTag(memberFromTag)

		if err != nil {
			return &Member{}, fmt.Errorf("Could not find member %s", memberFromTag)
		}
	}

	return memberFrom, nil
}

func getCommandParams(r *http.Request) (string, string, string, error) {
	var text = r.PostFormValue("text")

	textParts := strings.SplitN(text, " ", 3)
	if len(textParts) < 2 {
		return "", "", "", fmt.Errorf("Invalid number of parameters.")
	}

	command, target, extra := textParts[0], textParts[1], ""
	if len(textParts) == 3 {
		extra = textParts[2]
	}

	return command, target, extra, nil
}
