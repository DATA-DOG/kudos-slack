package main

import (
	"fmt"
	"html/template"
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

	r.Header.Set("Content-Type", "text/html")

	var page = `
	<!DOCTYPE html>
	<html lang="en">
	  <head>
	    <meta charset="utf-8">
	    <meta http-equiv="X-UA-Compatible" content="IE=edge">
	    <meta name="viewport" content="width=device-width, initial-scale=1">
	    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
	</head>
	<body>
	{{range .}}
		<blockquote>
			<p><span class="text-muted">@{{.MemberTo.Name}}:</span> {{.Kudo}}</p>
			<footer>{{.MemberFrom.RealName}} <cite title="Likes">({{.LikeCount}} likes)</cite></footer>
		</blockquote>
	{{end}}
	</body>
	</html>
	`
	tmpl := template.New("page")
	var err error
	if tmpl, err = tmpl.Parse(page); err != nil {
		fmt.Println(err)
	}

	tmpl.Execute(w, kudos)
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
