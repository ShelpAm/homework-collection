package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"
)

func ListAssignments(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	list := make([]Assignment, 0, len(assignments))
	for _, a := range assignments {
		a.BeginTime = a.BeginTime.Truncate(time.Second)
		a.EndTime = a.EndTime.Truncate(time.Second)
		list = append(list, a)
	}

	sort.Slice(list, func(i, j int) bool {
		if !list[i].BeginTime.Equal(list[j].BeginTime) {
			return list[i].BeginTime.Before(list[j].BeginTime)
		}
		if !list[i].EndTime.Equal(list[j].EndTime) {
			return list[i].EndTime.Before(list[j].EndTime)
		}
		return list[i].Name < list[j].Name
	})

	json.NewEncoder(w).Encode(list)
}
