package main

import (
	"net/http"
	"time"
	"html/template"
	"gopkg.in/mgo.v2"
	"log"
	"multiadv"
)

type CalendarDay struct {
	Day     int
	Entries []multiadv.Entry
}

type IndexModel struct {
	CalendarRows [][]CalendarDay
}

func filterEntriesOfDay(entries []multiadv.Entry, day int) []multiadv.Entry {
	res := make([]multiadv.Entry, 0)
	for _, entry := range entries {
		if entry.Day == day {
			res = append(res, entry)
		}
	}
	return res
}

func indexHandler(entryDAO multiadv.EntryDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		year := now.Year()
		month := now.Month()

		entries, err := entryDAO.FetchEntriesOfMonth(year, int(month))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		begin := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		indexModel := IndexModel{
			CalendarRows: [][]CalendarDay{make([]CalendarDay, begin.Weekday())},
		}
		for day := begin; day.Month() == month; day = day.AddDate(0, 0, 1) {
			if day.Day() != 1 && day.Weekday() == time.Sunday {
				indexModel.CalendarRows = append(indexModel.CalendarRows, make([]CalendarDay, 0))
			}
			i := len(indexModel.CalendarRows) - 1
			indexModel.CalendarRows[i] = append(indexModel.CalendarRows[i], CalendarDay{
				Day:     day.Day(),
				Entries: filterEntriesOfDay(entries, day.Day()),
			})
		}

		tmpl, err := template.ParseFiles("web/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, &indexModel)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/"+r.URL.String())
}

func main() {
	mongo, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Panicf("Cannot connect mongo\n")
	}
	entryDAO := multiadv.NewEntryDAO(mongo.DB("multiadv"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler(entryDAO))
	mux.HandleFunc("/static/", staticHandler)
	http.ListenAndServe("127.0.0.1:25252", mux)
}
