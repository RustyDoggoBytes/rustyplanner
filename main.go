package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	//go:embed static/*
	embeddedFiles embed.FS
	//go:embed schema.sql
	ddl string
)

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

type PageData struct {
	WeekStart    time.Time
	WeekEnd      time.Time
	PreviousWeek time.Time
	NextWeek     time.Time
	Meals        []MealPlan
	FormData     map[string][]string
}

var userID int64 = 1

func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func FormatMonthDay(date time.Time) string {
	return date.Format("01-02")
}

func main() {
	ctx := context.Background()
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	db, err := sql.Open("sqlite3", "data/planner.db")
	if err != nil {
		log.Fatal(err)
	}
	repository, err := NewRepository(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(getFileSystem())).ServeHTTP(w, r)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		date := time.Now()

		startDay := r.URL.Query().Get("start-date")
		if startDay != "" {
			date, err = time.Parse("2006-01-02", startDay)
			if err != nil {
				slog.Error("invalid date", "date", startDay)
			}
		}

		monday := getMondayOfCurrentWeek(date)

		slog.Info("received", "date", startDay, "monday", monday)
		sunday := monday.AddDate(0, 0, 6)
		meals, err := repository.GetMealPlanByDate(userID, monday, sunday)
		if err != nil {
			slog.Error("failed to retrieve meals", "start-date", monday, "end-date", monday, "error", err)
		}

		component := Index(PageData{
			WeekStart:    monday,
			WeekEnd:      sunday,
			PreviousWeek: monday.AddDate(0, 0, -7),
			NextWeek:     monday.AddDate(0, 0, 7),
			Meals:        meals,
		})

		component.Render(r.Context(), w)
	})

	mux.HandleFunc("POST /meal-plan/{date}", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("failed to read meals form", err)
		}

		date, err := time.Parse("2006-01-02", r.PathValue("date"))
		if err != nil {
			log.Fatal(err)
		}
		meal := MealPlan{
			Date:      date,
			Breakfast: r.FormValue("breakfast"),
			Snack1:    r.FormValue("snack1"),
			Snack2:    r.FormValue("snack2"),
			Lunch:     r.FormValue("lunch"),
			Dinner:    r.FormValue("dinner"),
		}

		err = repository.UpdateMealPlan(userID, meal)
		component := MealPlanCardForm(meal, err == nil, err)
		component.Render(r.Context(), w)
	})

	address := fmt.Sprintf("%s:8080", host)
	slog.Info("running server", "address", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
