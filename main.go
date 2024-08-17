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
	"rustydoggobytes/planner/components"
	"rustydoggobytes/planner/db"
	"rustydoggobytes/planner/middlewares"
	"rustydoggobytes/planner/utils"
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

var userID int64 = 1

func main() {
	ctx := context.Background()
	host := utils.GetEnv("HOST", "localhost")

	sqlite3, err := sql.Open("sqlite3", "data/planner.db")
	if err != nil {
		log.Fatal(err)
	}
	repository, err := db.NewRepository(ctx, sqlite3, ddl)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(getFileSystem())).ServeHTTP(w, r)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/groceries", http.StatusFound)
	})

	mux.HandleFunc("GET /meal-plans", func(w http.ResponseWriter, r *http.Request) {
		date := time.Now()

		startDay := r.URL.Query().Get("start-date")
		if startDay != "" {
			date, err = time.Parse("2006-01-02", startDay)
			if err != nil {
				slog.Error("invalid date", "date", startDay)
			}
		}

		monday := utils.GetMondayOfCurrentWeek(date)
		sunday := monday.AddDate(0, 0, 6)
		meals, err := repository.GetMealPlanByDate(userID, monday, sunday)
		if err != nil {
			slog.Error("failed to retrieve meals", "start-date", monday, "end-date", monday, "error", err)
		}

		component := components.MealPage(components.PageData{
			WeekStart:    monday,
			WeekEnd:      sunday,
			PreviousWeek: monday.AddDate(0, 0, -7),
			NextWeek:     monday.AddDate(0, 0, 7),
			Meals:        meals,
		})

		component.Render(r.Context(), w)
	})

	mux.HandleFunc("POST /meal-plans/{date}", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("failed to read meals form", err)
		}

		date, err := time.Parse("2006-01-02", r.PathValue("date"))
		if err != nil {
			log.Fatal(err)
		}
		meal := db.MealPlan{
			Date:      date,
			Breakfast: r.FormValue("breakfast"),
			Snack1:    r.FormValue("snack1"),
			Snack2:    r.FormValue("snack2"),
			Lunch:     r.FormValue("lunch"),
			Dinner:    r.FormValue("dinner"),
		}

		err = repository.UpdateMealPlan(userID, meal)
		component := components.MealPlanCardForm(meal, err == nil, err)
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("GET /groceries", func(w http.ResponseWriter, r *http.Request) {
		items, err := repository.ListGroceryItems(userID)
		if err != nil {
			slog.Error("failed to list groceries", "user_id", userID)
		}

		component := components.GroceryList(items)
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("POST /groceries", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("failed to read groceries form", err)
		}
		item, err := repository.CreateGroceryItem(userID, r.FormValue("name"))
		if err != nil {
			slog.Error("failed to create grocery", err)
		}

		component := components.GroceryListItem(*item)
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("DELETE /groceries/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		err := repository.DeleteGroceryItem(userID, id)
		if err != nil {
			slog.Error("failed to delete item", "id", id, err)
		}
	})

	mux.HandleFunc("PUT /groceries/{id}/toggle", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		item, err := repository.ToggleGroceryItem(userID, id)
		if err != nil {
			slog.Error("failed to toggle item", "id", id, err)
		}

		component := components.GroceryListItem(*item)
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("GET /chores", func(w http.ResponseWriter, r *http.Request) {
		chores := []string{
			"Take-out trash",
		}
		component := components.ChoresPage(chores)
		component.Render(r.Context(), w)
	})

	address := fmt.Sprintf("%s:8080", host)
	slog.Info("running server", "address", address)

	protectedMux := middlewares.BasicAuthMiddleware(
		mux,
		utils.GetEnv("AUTH_USER", "rusty"),
		utils.GetEnv("AUTH_PASSWORD", "doggo"),
	)

	log.Fatal(http.ListenAndServe(address, protectedMux))
}
