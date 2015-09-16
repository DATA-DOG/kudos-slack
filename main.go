package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type kudo struct {
	Kudo       string
	MemberTo   *Member
	MemberFrom *Member
	Likes      like
}

type like struct {
	Count   int
	Members *[]Member
}

var kudos []kudo

func main() {
	readConfig()
	loadUsers()

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/user/:id", getUserByID)

	router.POST("/kudo", handleKudoCmd)

	fmt.Print("Listening on port ", config.Port, "...")
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", config.Port), router))
}

func handleKudoCmd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var text = r.PostFormValue("text")
	var memberFromTag = r.PostFormValue("user_name")

	textParts := strings.SplitN(text, " ", 3)
	if len(textParts) != 3 {
		printKudoUsage(w)
		return
	}
	command, target, extra := textParts[0], textParts[1], textParts[2]

	var memberFrom, err = findMemberByTag(memberFromTag)
	if err != nil {
		loadUsers()
		memberFrom, err = findMemberByTag(memberFromTag)
	}

	switch command {
	case "to":
		member, err := findMemberByTag(strings.TrimLeft(target, "@"))
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		kudos = append(kudos, kudo{extra, member, memberFrom, like{}})
		notifyUser("New kudo from <@"+memberFrom.ID+">!\n"+extra, *member)
		fmt.Fprint(w, "Kudo has been registered!")
		break
	case "like":
		member, err := findMemberByTag(strings.TrimLeft(target, "@"))
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		for i := len(kudos) - 1; i >= 0; i-- {
			kudo := kudos[i]
			if kudo.MemberTo.ID == member.ID {
				kudo.Likes.Count++
				fmt.Fprint(w, "Got that!")
			}
		}
	default:
		printKudoUsage(w)
	}
}

func printKudoUsage(w http.ResponseWriter) {
	fmt.Fprint(w, "incorrect usage, baby!")
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
