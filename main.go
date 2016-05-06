package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"path/filepath"

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
var templates map[string]*template.Template

func main() {
	readConfig()
	loadUsers()
	loadDatabase()
	loadTemplates()

	router := httprouter.New()
	router.GET("/", loadKudosPage)
	router.GET("/calendar", loadCalendarPage)
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

func loadTemplates() {
    log.Println("Loading templates...")
    if templates == nil {
        templates = make(map[string]*template.Template)
    }

    templatesDir := config.TemplatesPath
    layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
    if err != nil {
        log.Fatal(err)
    }

    includes, err := filepath.Glob(templatesDir + "/includes/*.tmpl")
    if err != nil {
        log.Fatal(err)
    }

    // Generate our templates map from our layouts/ and includes/ directories
    for _, layout := range layouts {
        files := append(includes, layout)
        templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
    }
}
