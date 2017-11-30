package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type event struct {
	Happening bool
	Date      string
	EndDate   string
	Event     *calendar.Event
	Today     bool
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, authConfig *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(config.CalendarCredentials)
	if err != nil {
		tok = getTokenFromWeb(authConfig)
		saveToken(config.CalendarCredentials, tok)
	}
	return authConfig.Client(ctx, tok)
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getSheetClient(ctx context.Context, authConfig *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(config.AnniversariesSheetCredentials)
	if err != nil {
		tok = getTokenFromWeb(authConfig)
		saveToken(config.AnniversariesSheetCredentials, tok)
	}
	return authConfig.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getEvents() []event {
	ctx := context.Background()

	b, err := ioutil.ReadFile(config.CalendarSecret)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	calendarConfig, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, calendarConfig)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(config.CalendarID).ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
	}

	var compiledEvents []event

	var currentDate = time.Now()

	for _, i := range events.Items {
		var date string
		var happening bool
		var isToday bool
		// If the DateTime is an empty string the Event is an all-day Event.
		// So only Date is available.
		if i.Start.DateTime != "" {
			startDate, _ := time.Parse(time.RFC3339, i.Start.DateTime)
			endDate, _ := time.Parse(time.RFC3339, i.End.DateTime)

			var startDay = startDate.Format("02")
			isToday = currentDate.Format("02") == startDay || currentDate.After(startDate)
			happening = isToday && startDate.Before(currentDate) && endDate.After(currentDate)

			if startDay != endDate.Format("02") || !isToday {
				var dateStart = startDate.Format("01-02")
				var dateEnd = endDate.Format("01-02")

				if dateStart == dateEnd {
					date = dateStart + " " + startDate.Format("15:04") + "-" + endDate.Format("15:04")
				} else {
					date = startDate.Format("01-02 15:04") + "-" + endDate.Format("01-02 15:04")
				}
			} else {
				date = startDate.Format("15:04") + "-" + endDate.Format("15:04")
			}
		} else {
			startDate, _ := time.Parse("2006-01-02", i.Start.Date)
			date = startDate.Format("01-02") + " (all day)"

			if currentDate.After(startDate) {
				isToday = true
			} else {
				isToday = i.Start.Date == currentDate.Format("2006-01-02")
			}
		}

		compiledEvents = append(compiledEvents, event{Happening: happening, Date: date, Event: i, Today: isToday})
	}

	return compiledEvents
}

func getCalendarEvents() []event {
	ctx := context.Background()

	b, err := ioutil.ReadFile(config.CalendarSecret)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	calendarConfig, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, calendarConfig)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(config.FullCalendarID).ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(30).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
	}

	var compiledEvents []event
	for _, i := range events.Items {
		var dateStart, dateEnd string
		// If the DateTime is an empty string the Event is an all-day Event.
		// So only Date is available.
		if i.Start.DateTime != "" {
			startDate, _ := time.Parse(time.RFC3339, i.Start.DateTime)
			endDate, _ := time.Parse(time.RFC3339, i.End.DateTime)
			dateStart = startDate.Format("2006-01-02")
			dateEnd = endDate.Format("2006-01-02")
		} else {
			startDate, _ := time.Parse("2006-01-02", i.Start.Date)
			endDate, _ := time.Parse("2006-01-02", i.End.Date)
			dateStart = startDate.Format("2006-01-02")
			dateEnd = endDate.Format("2006-01-02")
		}
		compiledEvents = append(compiledEvents, event{Date: dateStart, EndDate: dateEnd, Event: i})
	}
	return compiledEvents
}
