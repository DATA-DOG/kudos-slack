package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Kudo main kudos
type Kudo struct {
	ID         int64
	Text       string
	Original   string
	MemberFrom Member
	Recipients []Member
	Value      int
	Color      string
	Date       time.Time
}

type pageView struct {
	Kudos  []kudoView
	Events []event
}

type kudoView struct {
	Item Kudo
	Text []string
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

	log.Println("Listening on port", config.Port, "...")
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", config.Port), router))
}

func handleKudoCmd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	memberFrom, err := getMemberFrom(r)
	if err != nil {
		fmt.Fprint(w, "Invalid user provided")
		return
	}

	command, text, err := getCommandParams(r)
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
		handleNewKudoCommand(w, memberFrom, text, value)
	default:
		printKudoUsage(w)
	}
}

func printKudoUsage(w http.ResponseWriter) {
	fmt.Fprint(w, "New kudo: `/kudos to @user1 [@user2, ...] message`")
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	r.Header.Set("Content-Type", "text/html")

	var page = `
	<!DOCTYPE html>
	<html>
	<head>
  	<meta charset="utf-8">
		<meta http-equiv="refresh" content="300">
	  <meta http-equiv="X-UA-Compatible" content="IE=edge">
	  <meta name="viewport" content="width=device-width, initial-scale=1">
	  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
	  <link rel="stylesheet" href="/asset/c.css">
	  <link href='https://fonts.googleapis.com/css?family=Patrick+Hand|Droid+Sans' rel='stylesheet' type='text/css'>
		<script type="text/javascript" src="/asset/snowstorm-min.js"></script>
		<script type="text/javascript">
			snowStorm.snowColor = '#FBFDFF';
			snowStorm.flakesMaxActive = 15;
			snowStorm.flakesMax = 15;
			snowStorm.animationInterval = 60;
			snowStorm.followMouse = false;
			snowStorm.useMeltEffect = false;

		</script>
	</head>
	<body>
	<img src="/asset/snowman.png" class="snowman">
	<div class="row">
  <div class="col-xs-7">
	<div class="notes">
	{{range .Kudos}}
		<div class="note note-{{.Item.Color}}">
			<div class="pin"></div>
			{{if eq .Item.Value -1}}
		  	<div class="sad"></div>
			{{end}}
			<p>{{range .Item.Recipients}}@{{.Name}}, {{end}}<br>{{range .Text}}{{.}}<br>{{end}}</p>
			<span>@{{.Item.MemberFrom.Name}}</span>
		</div>
	{{end}}
	</div>
	</div>
<div class="col-xs-5">
	<div class="calendar">
	<div class="pin pin-left"></div>
      <div class="pin pin-right"></div>
      <h1>Happening today</h1>
      <ul>
				{{range .Events}}
			  	{{if .Today}}
			  		<li{{if .Happening}} class="active"{{end}}><span class="text-muted">{{.Date}}</span> {{.Event.Summary}}</li>
			  	{{end}}
				{{end}}
      </ul>

			<h1>Upcoming events</h1>
      <ul>
				{{range .Events}}
			  	{{if not .Today}}
			  		<li><span class="text-muted">{{.Date}}</span> {{.Event.Summary}}</li>
			  	{{end}}
				{{end}}
      </ul>
	</div>
</div>
	</body>
	</html>
	`
	tmpl := template.New("page")
	var err error
	if tmpl, err = tmpl.Parse(page); err != nil {
		fmt.Println(err)
	}

	var viewKudos []kudoView
	for i := 0; i < 9 && i < len(kudos); i++ {
		view := kudoView{Item: kudos[i], Text: strings.Split(kudos[i].Text, "\n")}
		viewKudos = append(viewKudos, view)
	}
	pageData := pageView{Kudos: viewKudos, Events: getEvents()}

	tmpl.Execute(w, pageData)
}

func getMemberFrom(r *http.Request) (Member, error) {
	memberFromTag := r.PostFormValue("user_name")

	var memberFrom, err = findMemberByTag(memberFromTag)
	if err != nil {
		loadUsers()
		memberFrom, err = findMemberByTag(memberFromTag)

		if err != nil {
			return Member{}, fmt.Errorf("Could not find member %s", memberFromTag)
		}
	}

	return memberFrom, nil
}

func getCommandParams(r *http.Request) (string, string, error) {
	var text = r.PostFormValue("text")

	textParts := strings.SplitN(text, " ", 2)
	if len(textParts) < 1 {
		return "", "", fmt.Errorf("Invalid number of parameters.")
	}

	command, text := textParts[0], ""
	if len(textParts) == 2 {
		text = textParts[1]
	}

	return command, text, nil
}

func randomColor() string {
	colors := []string{"yellow", "pink", "green", "red", "orange", ""}
	return colors[rand.Intn(len(colors))]
}
