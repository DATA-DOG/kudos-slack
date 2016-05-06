package main

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/julienschmidt/httprouter"
)

func loadKudosPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  pageName := "kudos.tmpl"
	tmpl, ok := templates[pageName]
	if !ok {
			fmt.Errorf("The template %s does not exist.", pageName)
	}

	var viewKudos []kudoView
	for i := 0; i < 9 && i < len(kudos); i++ {
		view := kudoView{Item: kudos[i], Text: strings.Split(kudos[i].Text, "\n")}
		viewKudos = append(viewKudos, view)
	}
	pageData := pageView{Kudos: viewKudos}

	r.Header.Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "base.tmpl", pageData)
}
