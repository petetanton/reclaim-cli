package main

import (
	"encoding/json"
	"io"
	"time"

	"github.com/petetanton/reclaim-cli/pkg/reclaim"
)

// taskOutput is the machine-readable projection of a task written to stdout under
// --json. It is a deliberate subset rather than the whole upstream shape, so
// callers are not coupled to every field the API happens to return. logrus
// diagnostics go to stderr, which keeps this side of the output pipeable.
type taskOutput struct {
	Id                 int        `json:"id"`
	Title              string     `json:"title"`
	Notes              string     `json:"notes,omitempty"`
	Status             string     `json:"status"`
	Priority           string     `json:"priority"`
	Due                *time.Time `json:"due"`
	SnoozeUntil        *time.Time `json:"snoozeUntil"`
	TimeChunksRequired int        `json:"timeChunksRequired"`
	MinChunkSize       int        `json:"minChunkSize"`
	MaxChunkSize       int        `json:"maxChunkSize"`
	EventCategory      string     `json:"eventCategory"`
}

func newTaskOutput(task *reclaim.Task) taskOutput {
	return taskOutput{
		Id:                 task.Id,
		Title:              task.Title,
		Notes:              task.Notes,
		Status:             task.Status,
		Priority:           task.Priority,
		Due:                task.Due,
		SnoozeUntil:        task.SnoozeUntil,
		TimeChunksRequired: task.TimeChunksRequired,
		MinChunkSize:       task.MinChunkSize,
		MaxChunkSize:       task.MaxChunkSize,
		EventCategory:      task.EventCategory,
	}
}

// writeTaskJSON emits a task as a single line of JSON, so that
// `reclaim create --json | jq .id` works without a second API round-trip.
func writeTaskJSON(w io.Writer, task *reclaim.Task) error {
	return json.NewEncoder(w).Encode(newTaskOutput(task))
}
