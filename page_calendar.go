package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type calendarPageView struct {
	CalendarEvents []event
	Events []event
}
func loadCalendarPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  pageName := "calendar.tmpl"
	tmpl, ok := templates[pageName]
	if !ok {
			fmt.Errorf("The template %s does not exist.", pageName)
	}

	pageData := calendarPageView{CalendarEvents: getCalendarEvents(), Events: getEvents()}

	r.Header.Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "base.tmpl", pageData)
}
