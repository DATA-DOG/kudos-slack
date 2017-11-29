package main

import (
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"github.com/Iwark/spreadsheet"
	"regexp"
	"strings"
	"time"
)

type EmployeeAnniversary struct {
	Name       string
	Date       string
	WorksYears int
	Color      string
}

func getWorkAnniversaries() []EmployeeAnniversary {
	ctx := context.Background()

	b, err := ioutil.ReadFile(config.CalendarSecret)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	calendarConfig, err := google.ConfigFromJSON(b, spreadsheet.Scope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getSheetClient(ctx, calendarConfig)
	service := &spreadsheet.Service{Client: client}

	sheets, _ := service.Get(config.AnniversariesSheetId)

	ws1, _ := sheets.Get(0)

	currentMonth := int(time.Now().Month()) - 1

	employees := []EmployeeAnniversary{}
	for i, row := range ws1.Rows {
		if i == 0 {
			continue
		}
		rowData := row[currentMonth]
		if rowData == nil {
			break
		}

		employees = append(employees, parseEmployee(rowData))
	}

	return employees
}

func parseEmployee(data *spreadsheet.Cell) EmployeeAnniversary {
	r, _ := regexp.Compile("([0-9]{4}-[0-9]{2}-[0-9]{2})")
	rName, _ := regexp.Compile(`[0-9\(\)\-]+`)

	employeeWorkYear, _ := time.Parse("2006-01-02", r.FindString(data.Content))
	worksYears := time.Now().Year() - employeeWorkYear.Year()

	return EmployeeAnniversary{
		Color:      randomColor(),
		Date:       r.FindString(data.Content),
		Name:       strings.TrimSpace(rName.ReplaceAllString(data.Content, "")),
		WorksYears: worksYears,
	}
}
