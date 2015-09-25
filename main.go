package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
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
	Value      int
	Color      string
}

var kudos []Kudo

func main() {
	readConfig()
	loadUsers()
	loadDatabase()

	router := httprouter.New()
	router.GET("/", index)
	router.POST("/kudo", handleKudoCmd)
	router.POST("/boo", handleKudoCmd)

	router.ServeFiles("/asset/*filepath", http.Dir(config.AssetPath))

	fmt.Print("Listening on port ", config.Port, "...")
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", config.Port), router))
}

func handleKudoCmd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// accessToken := r.PostFormValue("token")
	//
	// if accessToken != config.SlackCommandToken {
	// 	fmt.Fprint(w, "Invalid access token")
	// 	return
	// }

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
		value := 1
		if r.PostFormValue("command") == "/boo" {
			value = -1
		}
		handleNewKudoCommand(w, memberFrom, target, extra, value)
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
	<html>
	<head>
  	<meta charset="utf-8">
	  <meta http-equiv="X-UA-Compatible" content="IE=edge">
	  <meta name="viewport" content="width=device-width, initial-scale=1">
	  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
	  <link rel="stylesheet" href="/asset/c.css">
	  <link href='https://fonts.googleapis.com/css?family=Patrick+Hand' rel='stylesheet' type='text/css'>
	</head>
	<body>
	<div class="notes">
	{{range .}}
		<div class="note note-{{.Color}}">
			<div class="pin"></div>
			{{if eq .Value -1}}
		  	<div class="sad"></div>
			{{end}}
			<p>@{{.MemberTo.Name}} {{.Value}},<br>{{.Kudo}}</p>
			<span>@{{.MemberFrom.Name}}</span>
		</div>
	{{end}}
	</div>
	</body>
	</html>
	`
	tmpl := template.New("page")
	var err error
	if tmpl, err = tmpl.Parse(page); err != nil {
		fmt.Println(err)
	}

	tmpl.Execute(w, reverseKudos())
}

func reverseKudos() <-chan Kudo {
	ch := make(chan Kudo)
	go func() {
		for i := len(kudos) - 1; i > 0; i-- {
			ch <- kudos[i]
		}
		close(ch)
	}()
	return ch
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

func randomColor() string {
	colors := []string{"yellow", "pink", "green", "red", "orange", ""}
	return colors[rand.Intn(len(colors))]
}
