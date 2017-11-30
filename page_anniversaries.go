package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"time"
)

type anniversariesPageView struct {
	Employees    []EmployeeAnniversary
	CurrentMonth string
}

func loadAnniversariesPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pageName := "anniversaries.tmpl"
	tmpl, ok := templates[pageName]

	if !ok {
		fmt.Errorf("The template %s does not exist.", pageName)
	}

	pageData := anniversariesPageView{
		Employees:    getWorkAnniversaries(),
		CurrentMonth: time.Now().Month().String(),
	}

	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "base.tmpl", pageData)
}
