package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func loadStatisticsPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  pageName := "statistics.tmpl"
	tmpl, ok := templates[pageName]
	if !ok {
			fmt.Errorf("The template %s does not exist.", pageName)
	}

	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "base.tmpl", getStatistics())
}
