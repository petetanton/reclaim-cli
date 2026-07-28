package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/petetanton/reclaim-cli/pkg/reclaim"
)

func TestWriteTaskJSON(t *testing.T) {
	snoozeUntil := time.Date(2026, 8, 4, 9, 15, 0, 0, time.UTC)
	task := &reclaim.Task{
		Id:                 13354745,
		Title:              `Review the "urgent" dashboard`,
		Status:             "NEW",
		Priority:           "P3",
		SnoozeUntil:        &snoozeUntil,
		TimeChunksRequired: 2,
		MinChunkSize:       2,
		MaxChunkSize:       16,
		EventCategory:      "WORK",
	}

	var out bytes.Buffer
	require.NoError(t, writeTaskJSON(&out, task))

	// One line, so the output composes with line-oriented tooling.
	assert.Equal(t, 1, strings.Count(out.String(), "\n"))

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &decoded))

	assert.Equal(t, float64(13354745), decoded["id"])
	assert.Equal(t, `Review the "urgent" dashboard`, decoded["title"])
	assert.Equal(t, "2026-08-04T09:15:00Z", decoded["snoozeUntil"])
	assert.Equal(t, float64(2), decoded["timeChunksRequired"])

	// due is always present so callers can distinguish "no due date" from a field
	// this CLI simply does not report.
	assert.Contains(t, decoded, "due")
	assert.Nil(t, decoded["due"])

	// Notes are omitted when empty rather than reported as "".
	assert.NotContains(t, decoded, "notes")
}

func TestWriteTaskJSON_includesDueWhenSet(t *testing.T) {
	due := time.Date(2026, 8, 14, 17, 30, 0, 0, time.UTC)
	task := &reclaim.Task{Id: 1, Due: &due, Notes: "some notes"}

	var out bytes.Buffer
	require.NoError(t, writeTaskJSON(&out, task))

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &decoded))
	assert.Equal(t, "2026-08-14T17:30:00Z", decoded["due"])
	assert.Equal(t, "some notes", decoded["notes"])
}
