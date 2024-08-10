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
	FormData map[string][]string
}

type MealPlan struct {
	Day       string
	Breakfast string
	Snack1    string
	Lunch     string
	Snack2    string
	Dinner    string
}

var userID int64 = 1

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
		meals, _ := repository.GetMealPlanByDate(userID, time.Now())

		component := Index(meals)
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("POST /meal-plan", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("failed to read meals form", err)
		}
		meals := processWeeklyMealFromForm(r.Form)
		err = repository.UpdateMealsForDate(userID, time.Now(), meals)
		if err != nil {
			slog.Error("failed to update meals", err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	address := fmt.Sprintf("%s:8080", host)
	slog.Info("running server", "address", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
