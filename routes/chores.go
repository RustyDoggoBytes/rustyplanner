package routes

import "rustydoggobytes/planner/db"

type ChorePageData struct {
	Chores []db.Chores
	Error  string
}
