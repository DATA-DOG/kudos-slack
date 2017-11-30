package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type memberPageView struct {
	Member Member
	Kudos  []Kudo
}

func loadMemberPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageName := "member.tmpl"
	memberId := ps.ByName("id")
	tmpl, ok := templates[pageName]

	if !ok {
		fmt.Errorf("The template %s does not exist.", pageName)
	}

	member, _ := findMemberByID(memberId)

	pageData := memberPageView{
		Member: member,
		Kudos:  loadKudosReceivedByUser(memberId),
	}

	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "base.tmpl", pageData)
}
