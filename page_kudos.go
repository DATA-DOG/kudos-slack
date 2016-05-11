package main

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/julienschmidt/httprouter"
)

// Kudo main kudos
type KudosStats struct {
	Member    Member
	Pts				int
	Position  int
	Prc			  float32
	HasCrown  bool
}

type kudosView struct {
	Item Kudo
	Text []string
}

type kudosPageView struct {
	Kudos		  		[]kudosView
	KudosReceived []KudosStats
	KudosGave			[]KudosStats
}

func loadKudosPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  pageName := "kudos.tmpl"
	tmpl, ok := templates[pageName]
	if !ok {
			fmt.Errorf("The template %s does not exist.", pageName)
	}

	var viewKudos []kudosView
	for i := 0; i < 9 && i < len(kudos); i++ {
		view := kudosView{Item: kudos[i], Text: strings.Split(kudos[i].Text, "\n")}
		viewKudos = append(viewKudos, view)
	}

	pageData := kudosPageView{
		Kudos: viewKudos,
		KudosReceived: loadKudosReceivedList(),
		KudosGave: loadKudosGaveList()}

	r.Header.Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "base.tmpl", pageData)
}
