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
	  <meta name="viewport" content="width=1920, initial-scale=1">
		<meta name="mobile-web-app-capable" content="yes">
	  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
	  <link rel="stylesheet" href="/asset/c.css">
	  <link href='https://fonts.googleapis.com/css?family=Patrick+Hand|Droid+Sans' rel='stylesheet' type='text/css'>
		<style>
@-webkit-keyframes snow
  {
  0%{background-position:0 0,0 0,0 0}
  to{background-position:500px 1000px,400px 400px,300px 300px}}

  @keyframes snow{0%{background-position:0 0,0 0,0 0}
  to{background-position:500px 1000px,400px 400px,300px 300px}
  }

  </style>
	</head>
	<body class="win95">
	<!--<div style="
	    position: fixed;
	    width: 100%;
	    height: 100%;
	    -webkit-animation-name: snow;
	    animation-name: snow;
	    -webkit-animation-duration: 20s;
	    animation-duration: 20s;
	    -webkit-animation-timing-function: linear;
	    animation-timing-function: linear;
	    -webkit-animation-delay: 0;
	    animation-delay: 0;
	    -webkit-animation-iteration-count: infinite;
	    animation-iteration-count: infinite;
	    background-image: url(/asset/snowh.png),url(/asset/snow3q.png),url(/asset/snow2l.png);
	    transform: translateZ(0);
	    z-index: 10000000;
	    /* background-color: red; */
	    transform: translateZ(0);
	"></div>

	<img src="/asset/snowman.png" class="snowman">-->
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
