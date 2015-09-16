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
	Kudo       string
	MemberTo   *Member
	MemberFrom *Member
	Likes      like
}

type like struct {
	Count   int
	Members []Member
}

var kudos []Kudo

func main() {
	readConfig()
	loadDatabase()
	loadUsers()

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/user/:id", getUserByID)

	router.POST("/kudo", handleKudoCmd)

	fmt.Print("Listening on port ", config.Port, "...")
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", config.Port), router))
}

func handleKudoCmd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
		fmt.Fprint(w, kudo.MemberFrom.Name, ": ", kudo.MemberTo.Name, ", ", kudo.Kudo, kudo.Likes.Count, "\n")
	}
}

func findMemberByTag(tag string) (*Member, error) {
	for _, user := range users {
		if user.Name == tag {
			return &user, nil
		}
	}

	return &Member{}, fmt.Errorf("Member with tag %s could not be found!", tag)
}

func getUserByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//			w.Header().Set("Content-Type", "text/html")

	//ps.ByName("id")

	//	fmt.Fprintf(w, "hello, %s!\n", "GUEST")
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
