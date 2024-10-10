package reclaim

import "time"

type Task struct {
	Id                  int       `json:"id"`
	Title               string    `json:"title"`
	Notes               string    `json:"notes"`
	EventCategory       string    `json:"eventCategory"`
	EventSubType        string    `json:"eventSubType"`
	Status              string    `json:"status"`
	TimeChunksRequired  int       `json:"timeChunksRequired"`
	TimeChunksSpent     int       `json:"timeChunksSpent"`
	TimeChunksRemaining int       `json:"timeChunksRemaining"`
	MinChunkSize        int       `json:"minChunkSize"`
	MaxChunkSize        int       `json:"maxChunkSize"`
	AlwaysPrivate       bool      `json:"alwaysPrivate"`
	Deleted             bool      `json:"deleted"`
	Index               float64   `json:"index"`
	Due                 time.Time `json:"due"`
	Created             time.Time `json:"created"`
	Updated             time.Time `json:"updated"`
	Finished            time.Time `json:"finished"`
	Adjusted            bool      `json:"adjusted"`
	AtRisk              bool      `json:"atRisk"`
	TimeSchemeId        string    `json:"timeSchemeId"`
	Priority            string    `json:"priority"`
	OnDeck              bool      `json:"onDeck"`
	Deferred            bool      `json:"deferred"`
	SortKey             float64   `json:"sortKey"`
	TaskSource          struct {
		Type string `json:"type"`
	} `json:"taskSource"`
	ReadOnlyFields          []interface{} `json:"readOnlyFields"`
	RecurringAssignmentType string        `json:"recurringAssignmentType"`
	Type                    string        `json:"type"`
}
