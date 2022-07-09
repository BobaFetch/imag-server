package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// get all info for a department
// dept name, number, top 20 jobs by priority/schedule date
// stats: daily goal, jobs completed, parts completed
func getQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dept := mux.Vars(r)["dept"]

	var queue QueueData

	queue.Department = getDepartmentInfo(dept)
	queue.Parts = getPartList(dept)
	queue.Stats.Goal = getDailyGoal(dept) / 34
	queue.Stats.CompletedJobs = getCompletedJobCount(dept)
	queue.Stats.CompletedParts = getCompletedPartCount(dept)
	queue.Employees = getEmployeeDailyStats(dept)

	json.NewEncoder(w).Encode(queue)
}
