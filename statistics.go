package main

import (
    "io/ioutil"
    "log"
    "strconv"

    "golang.org/x/net/context"
  	"golang.org/x/oauth2/google"
    "github.com/Iwark/spreadsheet"
)

type BillableStatistics struct {
	Dates []string
  Ideal []float64
  Expected []float64
  Real []float64
}

type DeadlineStatistics struct {
	Project string
  Deadline string
  Left int64
  IsToday bool
}

type Statistics struct {
  BillableStatistics BillableStatistics
  DeadlineStatistics []DeadlineStatistics
}

func getStatistics() Statistics {
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

    sheets, _ := service.Get(config.StatisticsSheetId)

    var statistics Statistics
    var billableStatistics BillableStatistics
    ws1, _ := sheets.Get(0)
    for i, row := range ws1.Rows {
        if i == 0 {
          continue
        }
        billableStatistics.Dates = append(billableStatistics.Dates, row[0].Content)
        ideal, _ := strconv.ParseFloat(row[1].Content, 32)
        real, _ := strconv.ParseFloat(row[2].Content, 32)
        expected := ideal * 0.8
        billableStatistics.Ideal = append(billableStatistics.Ideal, ideal)
        billableStatistics.Expected = append(billableStatistics.Expected, expected)
        billableStatistics.Real = append(billableStatistics.Real, real)
    }
    statistics.BillableStatistics = billableStatistics

    var deadlineStatistics []DeadlineStatistics
    ws2, _ := sheets.Get(1)
    for i, row := range ws2.Rows {
        left, _ := strconv.ParseInt(row[2].Content, 10, 32)
        if i == 0 || left < 0 {
          continue
        }
        var deadline DeadlineStatistics
        deadline.Project = row[0].Content
        deadline.Deadline = row[1].Content
        deadline.IsToday = false
        deadline.Left = left
        if left == 0 {
          deadline.IsToday = true
        }
        deadlineStatistics = append(deadlineStatistics, deadline)
    }
    statistics.DeadlineStatistics = deadlineStatistics

    return statistics
}
