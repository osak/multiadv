package main

import (
	"net/http"
	"time"
	"html/template"
)

type CalendarDay struct {
	Day int
}

type IndexModel struct {
	CalendarRows [][]CalendarDay
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	year := now.Year()
	month := now.Month()

	begin := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	indexModel := IndexModel {
		CalendarRows: [][]CalendarDay{make([]CalendarDay, begin.Weekday())},
	}
	for day := begin; day.Month() == month; day = day.AddDate(0, 0, 1) {
		if day.Day() != 1 && day.Weekday() == time.Sunday {
			indexModel.CalendarRows = append(indexModel.CalendarRows, make([]CalendarDay, 0))
		}
		i := len(indexModel.CalendarRows) - 1
		indexModel.CalendarRows[i] = append(indexModel.CalendarRows[i], CalendarDay {
			Day: day.Day(),
		})
	}

	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, &indexModel)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/" + r.URL.String())
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/static/", staticHandler)
	http.ListenAndServe("127.0.0.1:25252", mux)
}
